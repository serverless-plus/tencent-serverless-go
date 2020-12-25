package core_test

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	faascontext "github.com/serverless-plus/tencent-serverless-go/context"
	"github.com/serverless-plus/tencent-serverless-go/core"
	"github.com/serverless-plus/tencent-serverless-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestAccessor tests", func() {
	Context("event conversion", func() {
		accessor := core.RequestAccessor{}
		basicRequest := getProxyRequest("/hello", "GET")
		It("Correctly converts a basic event", func() {
			httpReq, err := accessor.EventToRequestWithContext(context.Background(), basicRequest)
			Expect(err).To(BeNil())
			Expect("/hello").To(Equal(httpReq.URL.Path))
			Expect("/hello").To(Equal(httpReq.RequestURI))
			Expect("GET").To(Equal(httpReq.Method))
		})

		basicRequest = getProxyRequest("/hello", "get")
		It("Converts method to uppercase", func() {
			// calling old method to verify reverse compatibility
			httpReq, err := accessor.ProxyEventToHTTPRequest(basicRequest)
			Expect(err).To(BeNil())
			Expect("/hello").To(Equal(httpReq.URL.Path))
			Expect("/hello").To(Equal(httpReq.RequestURI))
			Expect("GET").To(Equal(httpReq.Method))
		})

		binaryBody := make([]byte, 256)
		_, err := rand.Read(binaryBody)
		if err != nil {
			Fail("Could not generate random binary body")
		}

		encodedBody := base64.StdEncoding.EncodeToString(binaryBody)

		binaryRequest := getProxyRequest("/hello", "POST")
		binaryRequest.Body = encodedBody
		binaryRequest.IsBase64Encoded = true

		It("Decodes a base64 encoded body", func() {
			httpReq, err := accessor.EventToRequestWithContext(context.Background(), binaryRequest)
			Expect(err).To(BeNil())
			Expect("/hello").To(Equal(httpReq.URL.Path))
			Expect("/hello").To(Equal(httpReq.RequestURI))
			Expect("POST").To(Equal(httpReq.Method))

			bodyBytes, err := ioutil.ReadAll(httpReq.Body)

			Expect(err).To(BeNil())
			Expect(len(binaryBody)).To(Equal(len(bodyBytes)))
			Expect(binaryBody).To(Equal(bodyBytes))
		})

		qsRequest := getProxyRequest("/hello", "GET")
		qsRequest.QueryString = map[string][]string{
			"hello": {"1", "2"},
		}
		It("Populates query string correctly", func() {
			httpReq, err := accessor.EventToRequestWithContext(context.Background(), qsRequest)
			Expect(err).To(BeNil())
			Expect("/hello").To(Equal(httpReq.URL.Path))
			Expect(httpReq.RequestURI).To(ContainSubstring("hello=1"))
			Expect("GET").To(Equal(httpReq.Method))

			query := httpReq.URL.Query()
			Expect(1).To(Equal(len(query)))
			Expect(query["hello"]).ToNot(BeNil())
			Expect(2).To(Equal(len(query["hello"])))
			Expect([]string{"1", "2"}).To(Equal(query["hello"]))
		})

		svhRequest := getProxyRequest("/hello", "GET")
		svhRequest.Headers = map[string]string{
			"hello": "1",
			"world": "2",
		}
		It("Populates single value headers correctly", func() {
			httpReq, err := accessor.EventToRequestWithContext(context.Background(), svhRequest)
			Expect(err).To(BeNil())
			Expect("/hello").To(Equal(httpReq.URL.Path))
			Expect("GET").To(Equal(httpReq.Method))

			headers := httpReq.Header
			Expect(2).To(Equal(len(headers)))

			for k, value := range headers {
				Expect(value[0]).To(Equal(svhRequest.Headers[strings.ToLower(k)]))
			}
		})

		basePathRequest := getProxyRequest("/app1/orders", "GET")

		It("Stips the base path correct", func() {
			accessor.StripBasePath("app1")
			httpReq, err := accessor.EventToRequestWithContext(context.Background(), basePathRequest)

			Expect(err).To(BeNil())
			Expect("/orders").To(Equal(httpReq.URL.Path))
			Expect("/orders").To(Equal(httpReq.RequestURI))
		})

		contextRequest := getProxyRequest("orders", "GET")
		contextRequest.Context = getRequestContext()

		It("Populates context header correctly", func() {
			// calling old method to verify reverse compatibility
			httpReq, err := accessor.ProxyEventToHTTPRequest(contextRequest)
			Expect(err).To(BeNil())
			Expect(2).To(Equal(len(httpReq.Header)))
			Expect(httpReq.Header.Get(core.APIGwContextHeader)).ToNot(BeNil())
		})
	})

	Context("StripBasePath tests", func() {
		accessor := core.RequestAccessor{}
		It("Adds prefix slash", func() {
			basePath := accessor.StripBasePath("app1")
			Expect("/app1").To(Equal(basePath))
		})

		It("Removes trailing slash", func() {
			basePath := accessor.StripBasePath("/app1/")
			Expect("/app1").To(Equal(basePath))
		})

		It("Ignores blank strings", func() {
			basePath := accessor.StripBasePath("  ")
			Expect("").To(Equal(basePath))
		})
	})

	Context("Retrieves API Gateway context", func() {
		It("Returns a correctly unmarshalled object", func() {
			contextRequest := getProxyRequest("orders", "GET")
			contextRequest.Context = getRequestContext()

			accessor := core.RequestAccessor{}
			// calling old method to verify reverse compatibility
			httpReq, err := accessor.ProxyEventToHTTPRequest(contextRequest)
			Expect(err).To(BeNil())

			headerContext, err := accessor.GetAPIGatewayContext(httpReq)
			Expect(err).To(BeNil())
			Expect(headerContext).ToNot(BeNil())
			Expect("x").To(Equal(headerContext.RequestID))
			proxyContext, ok := core.GetAPIGatewayContextFromContext(httpReq.Context())
			// should fail because using header proxy method
			Expect(ok).To(BeFalse())

			httpReq, err = accessor.EventToRequestWithContext(context.Background(), contextRequest)
			Expect(err).To(BeNil())
			proxyContext, ok = core.GetAPIGatewayContextFromContext(httpReq.Context())
			Expect(ok).To(BeTrue())
			Expect("x").To(Equal(proxyContext.RequestID))
			Expect("prod").To(Equal(proxyContext.Stage))
			runtimeContext, ok := core.GetRuntimeContextFromContext(httpReq.Context())
			Expect(ok).To(BeTrue())
			Expect(runtimeContext).To(BeNil())

			faasContext := faascontext.NewContext(context.Background(), &faascontext.FunctionContext{RequestID: "abc123"})
			httpReq, err = accessor.EventToRequestWithContext(faasContext, contextRequest)
			Expect(err).To(BeNil())

			headerContext, err = accessor.GetAPIGatewayContext(httpReq)
			// should fail as new context method doesn't populate headers
			Expect(err).ToNot(BeNil())
			proxyContext, ok = core.GetAPIGatewayContextFromContext(httpReq.Context())
			Expect(ok).To(BeTrue())
			Expect("x").To(Equal(proxyContext.RequestID))
			Expect("prod").To(Equal(proxyContext.Stage))
			runtimeContext, ok = core.GetRuntimeContextFromContext(httpReq.Context())
			Expect(ok).To(BeTrue())
			Expect(runtimeContext).ToNot(BeNil())
		})

		It("Populates the default hostname correctly", func() {

			basicRequest := getProxyRequest("orders", "GET")
			basicRequest.Context = getRequestContext()
			accessor := core.RequestAccessor{}
			httpReq, err := accessor.ProxyEventToHTTPRequest(basicRequest)
			Expect(err).To(BeNil())

			Expect(basicRequest.Context.SourceIP).To(Equal(httpReq.Host))
			Expect(basicRequest.Context.SourceIP).To(Equal(httpReq.URL.Host))
		})

		It("Uses a custom hostname", func() {
			myCustomHost := "http://my-custom-host.com"
			os.Setenv(core.CustomHostVariable, myCustomHost)
			basicRequest := getProxyRequest("orders", "GET")
			accessor := core.RequestAccessor{}
			httpReq, err := accessor.EventToRequestWithContext(context.Background(), basicRequest)
			Expect(err).To(BeNil())

			Expect(myCustomHost).To(Equal("http://" + httpReq.Host))
			Expect(myCustomHost).To(Equal("http://" + httpReq.URL.Host))
			os.Unsetenv(core.CustomHostVariable)
		})

		It("Strips terminating / from hostname", func() {
			myCustomHost := "http://my-custom-host.com"
			os.Setenv(core.CustomHostVariable, myCustomHost+"/")
			basicRequest := getProxyRequest("orders", "GET")
			accessor := core.RequestAccessor{}
			httpReq, err := accessor.EventToRequestWithContext(context.Background(), basicRequest)
			Expect(err).To(BeNil())

			Expect(myCustomHost).To(Equal("http://" + httpReq.Host))
			Expect(myCustomHost).To(Equal("http://" + httpReq.URL.Host))
			os.Unsetenv(core.CustomHostVariable)
		})
	})
})

func getProxyRequest(path string, method string) events.APIGatewayRequest {
	return events.APIGatewayRequest{
		Path:   path,
		Method: method,
	}
}

func getRequestContext() events.APIGatewayRequestContext {
	return events.APIGatewayRequestContext{
		RequestID: "x",
		Stage:     "prod",
		SourceIP:  "127.0.0.1",
	}
}

func getStageVariables() map[string]string {
	return map[string]string{
		"var1": "value1",
		"var2": "value2",
	}
}
