package chiadapter_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiadapter "github.com/serverless-plus/tencent-serverless-go/chi"
	"github.com/serverless-plus/tencent-serverless-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ChiLambda tests", func() {
	Context("Simple ping request", func() {
		It("Proxies the event correctly", func() {
			log.Println("Starting test")
			r := chi.NewRouter()
			r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				log.Println("Handler!!")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				buf, _ := json.Marshal(map[string]interface{}{
					"message": "pong",
				})
				w.Write(buf)
			})

			adapter := chiadapter.New(r)

			req := events.APIGatewayRequest{
				Path:   "/ping",
				Method: "GET",
			}

			resp, err := adapter.ProxyWithContext(context.Background(), req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))

			resp, err = adapter.Proxy(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
		})
	})
})
