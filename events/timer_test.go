package events_test

import (
	"encoding/json"

	"github.com/serverless-plus/tencent-serverless-go/events"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Events tests", func() {
	Context("timer event", func() {
		It("Correctly unmarshal timer event", func() {
			data := `{"Message": "", "Time": "2019-08-17T05:26:00Z", "TriggerName": "oneminute", "Type": "Timer"}`
			var event events.TimerEvent
			err := json.Unmarshal([]byte(data), &event)
			Expect(err).To(BeNil())
		})
	})
})
