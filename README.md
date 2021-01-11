## Tencent Serverless Go

[![Build Status](https://github.com/serverless-plus/tencent-serverless-go/workflows/Test/badge.svg?branch=master)](https://github.com/serverless-plus/tencent-serverless-go/actions?query=workflow:Test+branch:master)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/serverless-plus/tencent-serverless-go/gin?tab=doc)

## Getting started

The first step is to install the required dependencies

```bash
$ go get github.com/serverless-plus/tencent-serverless-go/events
$ go get github.com/serverless-plus/tencent-serverless-go/faas
$ go get github.com/serverless-plus/tencent-serverless-go/gin
```

```go
package main

import (
  "context"
  "fmt"

	"github.com/serverless-plus/tencent-serverless-go/events"
	"github.com/serverless-plus/tencent-serverless-go/faas"
	"github.com/serverless-plus/tencent-serverless-go/gin"
	"github.com/gin-gonic/gin"
)

var ginFaas *ginadapter.GinFaas

func init() {
	fmt.Printf("Gin start")
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
      "message": "Hello Serverless Gin",
      "query": c.Query("q"),
		})
	})

	ginFaas = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	var res, _ = ginFaas.ProxyWithContext(ctx, req)
  var apiRes = events.APIGatewayResponse{Body: res.Body, StatusCode: 200, Headers: res.Headers}
  return apiRes, nil
}

func main() {
  faas.Start(Handler)
}
```

## License

This library is licensed under the Apache 2.0 License.

Copyright 2020 Serverless Plus
