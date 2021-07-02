package chiadapter

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/serverless-plus/tencent-serverless-go/core"
	"github.com/serverless-plus/tencent-serverless-go/events"
)

// ChiFaas makes it easy to send API Gateway proxy events to a Chi
type ChiFaas struct {
	core.RequestAccessor

	chiEngine *chi.Mux
}

// New creates a new instance of the ChiFaas object.
// Receives an initialized *chi.Mux object - normally created with chi.Default().
// It returns the initialized instance of the ChiFaas object.
func New(chi *chi.Mux) *ChiFaas {
	return &ChiFaas{chiEngine: chi}
}

// Proxy receives an API Gateway proxy event, transforms it into an http.Request
// object, and sends it to the chi.Mux for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (g *ChiFaas) Proxy(req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	chiRequest, err := g.ProxyEventToHTTPRequest(req)
	return g.proxyInternal(chiRequest, err)
}

// ProxyWithContext receives context and an API Gateway proxy event,
// transforms them into an http.Request object, and sends it to the chi.Mux for routing.
// It returns a proxy response object generated from the http.ResponseWriter.
func (g *ChiFaas) ProxyWithContext(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	chiRequest, err := g.EventToRequestWithContext(ctx, req)
	return g.proxyInternal(chiRequest, err)
}

func (g *ChiFaas) proxyInternal(req *http.Request, err error) (events.APIGatewayResponse, error) {

	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriter()
	g.chiEngine.ServeHTTP(http.ResponseWriter(respWriter), req)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}
