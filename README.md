## Tencent Serverless Go

[TOC]

[![Build Status](https://github.com/serverless-plus/tencent-serverless-go/workflows/Test/badge.svg?branch=master)](https://github.com/serverless-plus/tencent-serverless-go/actions?query=workflow:Test+branch:master)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/serverless-plus/tencent-serverless-go/gin?tab=doc)


## Getting started (Gin)

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



## Getting started (Beego)

The first step is to install the required dependencies

```bash
$ go get github.com/serverless-plus/tencent-serverless-go/events
$ go get github.com/serverless-plus/tencent-serverless-go/faas
$ go get github.com/serverless-plus/tencent-serverless-go/beego
```

```go
package main

import (
	"github.com/beego/beego/v2/server/web"
	beegoadapter "github.com/serverless-plus/tencent-serverless-go/beego"
)

// Controller is controller for BeegoApp
type Controller struct {
	web.Controller
}

// Hello is Handler for "GET /" Route
func (ctrl *Controller) Hello() {
	ctrl.Ctx.Output.Body([]byte("Hello Serverless Beego"))
}

func main() {

	ctrl := &Controller{}

	web.Router("/", ctrl, "get:Hello")
	beegoadapter.Run(web.BeeApp)
}
```

the example can be found in `./example/beego`, after `make` in that directory, you can get an accessable url of the deployment, for example: 

[https://service-5tlgahl8-1256777886.gz.apigw.tencentcs.com/release/](https://service-5tlgahl8-1256777886.gz.apigw.tencentcs.com/release/)

### Migrate Beego application

you can also migrate your beego application with one step:

```go
// Replace:
// web.Run()
// To:
beegoadapter.Run(web.BeeApp)
```

## License

This library is licensed under the Apache 2.0 License.

Copyright 2020 Serverless Plus
