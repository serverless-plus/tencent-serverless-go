package ginadapter

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/serverless-plus/tencent-serverless-go/core"
	"github.com/serverless-plus/tencent-serverless-go/events"
)

// GinFaas makes it easy to send API Gateway proxy events to a Gin
type GinFaas struct {
	core.RequestAccessor

	ginEngine *gin.Engine
}

// New creates a new instance of the GinFaas object.
// Receives an initialized *gin.Engine object - normally created with gin.Default().
// It returns the initialized instance of the GinFaas object.
func New(gin *gin.Engine) *GinFaas {
	return &GinFaas{ginEngine: gin}
}

// Proxy receives an API Gateway proxy event, transforms it into an http.Request
// object, and sends it to the gin.Engine for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (g *GinFaas) Proxy(req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	ginRequest, err := g.ProxyEventToHTTPRequest(req)
	return g.proxyInternal(ginRequest, err)
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the gin.Engine for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (g *GinFaas) ProxyWithContext(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	ginRequest, err := g.EventToRequestWithContext(ctx, req)
	return g.proxyInternal(ginRequest, err)
}

func (g *GinFaas) proxyInternal(req *http.Request, err error) (events.APIGatewayResponse, error) {

	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriter()
	g.ginEngine.ServeHTTP(http.ResponseWriter(respWriter), req)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}
