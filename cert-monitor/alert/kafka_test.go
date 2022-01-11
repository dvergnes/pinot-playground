// Copyright (c) 2022 Denis Vergnes
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
