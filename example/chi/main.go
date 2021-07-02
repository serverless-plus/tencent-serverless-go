package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiadapter "github.com/serverless-plus/tencent-serverless-go/chi"
	"github.com/serverless-plus/tencent-serverless-go/events"
	"github.com/serverless-plus/tencent-serverless-go/faas"
)

var chiFaas *chiadapter.ChiFaas

func init() {
	fmt.Printf("Chi start")
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		buf, _ := json.Marshal(map[string]interface{}{
			"message": "Hello Serverless Chi",
			"query":   r.URL.Query().Get("q"),
		})
		w.Write(buf)
	})

	chiFaas = chiadapter.New(r)
}

// Handler serverless faas handler
func Handler(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	return chiFaas.ProxyWithContext(ctx, req)
}

func main() {
	faas.Start(Handler)
}
