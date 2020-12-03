package events_test

import (
	"encoding/json"

	"github.com/serverless-plus/tencent-serverless-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Events tests", func() {
	Context("cmq event", func() {
		It("Correctly unmarshal cmq event", func() {
			data := `{"Records": [{"CMQ": {"msgBody": "aaaaaaaaaa", "msgId": "13510798882111490", "msgTag": "aaaa, bbbbb, cccc", "publishTime": "2019-08-16T10:48:49Z", "requestId": "2758374289357404466", "subscriptionName": "ritchiechen-cmq-test", "topicName": "ritchiechen-cmq-test", "topicOwner": 123456, "type": "topic"}}]}`
			var event events.CMQEvent
			err := json.Unmarshal([]byte(data), &event)
			Expect(err).To(BeNil())
		})
	})
})

