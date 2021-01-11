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
