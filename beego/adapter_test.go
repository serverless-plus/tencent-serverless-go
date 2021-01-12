package beegoadapter_test

import (
	"context"
	"log"

	"github.com/beego/beego/v2/server/web"
	beegoadapter "github.com/serverless-plus/tencent-serverless-go/beego"
	"github.com/serverless-plus/tencent-serverless-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MainController struct {
	web.Controller
}

var _ = Describe("BeegoAdapter tests", func() {
	Context("GET /ping request", func() {
		app := web.BeeApp
		ctrl := &MainController{}
		log.Println("Starting test")
		web.Router("/ping", ctrl, "get:Ping")

		adapter := beegoadapter.New(app)
		req := events.APIGatewayRequest{
			Path:   "/ping",
			Method: "GET",
		}
		It("Proxy With Context", func() {

			resp, err := adapter.ProxyWithContext(context.Background(), req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Body).To(Equal(("Hello World!")))
		})

		It("Proxy Without Context", func() {
			resp, err := adapter.Proxy(req)

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Body).To(Equal(("Hello World!")))
		})
	})
})

func (ctrl *MainController) Ping() {
	ctrl.Ctx.Output.Body([]byte("Hello World!"))
}
