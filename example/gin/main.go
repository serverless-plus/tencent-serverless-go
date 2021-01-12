package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/serverless-plus/tencent-serverless-go/events"
	"github.com/serverless-plus/tencent-serverless-go/faas"
	ginadapter "github.com/serverless-plus/tencent-serverless-go/gin"
)

var ginFaas *ginadapter.GinFaas

func init() {
	fmt.Printf("Gin start")
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Serverless Gin",
			"query":   c.Query("q"),
		})
	})

	ginFaas = ginadapter.New(r)
}

// Handler serverless faas handler
func Handler(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	return ginFaas.ProxyWithContext(ctx, req)
}

func main() {
	faas.Start(Handler)
}
