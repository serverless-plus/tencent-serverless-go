package events_test

import (
	"encoding/json"

	"github.com/serverless-plus/tencent-serverless-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestAccessor tests", func() {
	Context("apigw event", func() {
		It("Correctly unmarshal apigw event", func() {
			data := `{"headerParameters": {}, "headers": {"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8", "accept-encoding": "gzip, deflate, br", "accept-language": "zh-CN,zh;q=0.9,en;q=0.8", "connection": "keep-alive", "endpoint-timeout": "15", "host": "service-xxx-123456.ap-shanghai.apigateway.myqcloud.com", "upgrade-insecure-requests": "1", "user-agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36", "x-anonymous-consumer": "true", "x-qualifier": "$LATEST"}, "httpMethod": "GET", "path": "/ritchiechen-apigw-test", "pathParameters": {}, "queryString": {"a": ["b", "c"], "x": "y", "i": true}, "queryStringParameters": {}, "requestContext": {"httpMethod": "ANY", "identity": {}, "path": "/ritchiechen-apigw-test", "serviceId": "service-xxx", "sourceIp": "8.8.8.8", "stage": "test"}}`
			var event events.APIGatewayRequest
			err := json.Unmarshal([]byte(data), &event)
			Expect(err).To(BeNil())
		})
	})
})
