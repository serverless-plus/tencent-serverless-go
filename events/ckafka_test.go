package events_test

import (
	"encoding/json"

	"github.com/serverless-plus/tencent-serverless-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Events tests", func() {
	Context("ckafka event", func() {
		It("Correctly unmarshal ckafka event", func() {
			data := `{"Records": [{"Ckafka": {"msgBody": "hello world", "msgKey": "ckafka-test-key", "offset": 7, "partition": 0, "topic": "ritchiechen-ckafka-test"}}]}`
			var event events.CkafkaEvent
			err := json.Unmarshal([]byte(data), &event)
			Expect(err).To(BeNil())
		})
	})
})
