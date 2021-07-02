## Tencent Serverless Go

[TOC]

[![Build Status](https://github.com/serverless-plus/tencent-serverless-go/workflows/Test/badge.svg?branch=master)](https://github.com/serverless-plus/tencent-serverless-go/actions?query=workflow:Test+branch:master)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/serverless-plus/tencent-serverless-go/gin?tab=doc)

## Usage

Support Web Framework:

- [gin](#gin)
- [beego](#beego)
- [chi](#chi)
### gin

[Example](./example/gin)

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

### beego

[Example](./example/beego)

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

### chi

[Example](./example/chi)

The first step is to install the required dependencies

```bash
$ go get github.com/serverless-plus/tencent-serverless-go/events
$ go get github.com/serverless-plus/tencent-serverless-go/faas
$ go get github.com/serverless-plus/tencent-serverless-go/chi
```

```go
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

```

## License

This library is licensed under the Apache 2.0 License.

Copyright 2020 Serverless Plus
