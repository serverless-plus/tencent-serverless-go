package beegoadapter

import (
	"context"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"github.com/serverless-plus/tencent-serverless-go/core"
	"github.com/serverless-plus/tencent-serverless-go/events"
	"github.com/serverless-plus/tencent-serverless-go/faas"
)

// BeegoAdapter inherit the RequestAccessor to proxy events
type BeegoAdapter struct {
	core.RequestAccessor
	beegoApp *web.HttpServer
}

// New create a BeegoAdapter Instance
func New(beego *web.HttpServer) *BeegoAdapter {
	return &BeegoAdapter{beegoApp: beego}
}

// Run a beego webApp
func Run(beego *web.HttpServer) *BeegoAdapter {
	adapter := New(beego)
	faas.Start(adapter.ProxyWithContext)
	return adapter
}

// Proxy convert APIGatewayRequest to HTTPRequest
func (b *BeegoAdapter) Proxy(req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	httpReq, err := b.ProxyEventToHTTPRequest(req)
	return b.proxyInternal(httpReq, err)
}

// ProxyWithContext convert APIGatewayRequest and Context to HTTPRequest
func (b *BeegoAdapter) ProxyWithContext(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	httpReq, err := b.EventToRequestWithContext(ctx, req)
	return b.proxyInternal(httpReq, err)
}

func (b *BeegoAdapter) proxyInternal(httpReq *http.Request, err error) (events.APIGatewayResponse, error) {
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("[BeegoAdapter] (proxyInternal) Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriter()
	b.beegoApp.Handlers.ServeHTTP(http.ResponseWriter(respWriter), httpReq)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("[BeegoAdapter] (proxyInternal) Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}
