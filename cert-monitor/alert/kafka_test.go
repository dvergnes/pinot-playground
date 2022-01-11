package alert_test

import (
	"errors"

	"github.com/dvergnes/pinot-playground/cert-monitor/alert"

	"github.com/Shopify/sarama/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka", func() {
	var (
		kafkaNotifier alert.Notifier
		producerMock  *mocks.SyncProducer

		err error
	)

	BeforeEach(func() {
		producerMock = mocks.NewSyncProducer(GinkgoT(), mocks.NewTestConfig())
		kafkaNotifier = alert.NewKafkaNotifier("topic", producerMock)
	})

	AfterEach(func() {
		Expect(producerMock.Close()).Should(Succeed())
	})

	Describe("Send", func() {
		var (
			a = alert.Alert{
				Level:   alert.Error,
				Message: "This is fine",
				ObjectRef: alert.ObjectRef{
					Name:      "cert",
					Namespace: "ns",
				},
				Source: "UT",
				When:   0,
			}
		)

		JustBeforeEach(func() {
			err = kafkaNotifier.Send(a)
		})

		When("message can be sent", func() {
			BeforeEach(func() {
				producerMock.ExpectSendMessageAndSucceed()
			})
			It("should not return any errors", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("producer failed to send message", func() {
			criticalError := errors.New("broker is down")
			BeforeEach(func() {
				producerMock.ExpectSendMessageAndFail(criticalError)
			})
			It("should not return any errors", func() {
				Expect(err).Should(MatchError("failed to deliver alert: broker is down"))
			})
		})

	})

	Describe("Close", func() {

		JustBeforeEach(func() {
			err = kafkaNotifier.Close()
		})

		When("producer close without errors", func() {
			It("should not return any errors", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

	})
})
