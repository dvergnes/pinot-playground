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
	"encoding/json"

	"github.com/dvergnes/pinot-playground/cert-monitor/alert"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alert", func() {

	Describe("Marshal JSON", func() {
		var (
			a alert.Alert
		)
		BeforeEach(func() {
			a = alert.Alert{
				Level:     alert.Warn,
				ObjectRef: alert.ObjectRef{
					Namespace: "ns",
					Name:      "cert",
				},
				Message:   "cert is about to expire",
				When:      1234567890,
				Source:    "host-123",
			}
		})
		It("should return a JSON object", func() {
			data, err:=json.Marshal(a)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(Equal([]byte(`{"level":"WARN","message":"cert is about to expire","objectRef":{"name":"cert","namespace":"ns"},"source":"host-123","when":1234567890}`)))
		})
	})
})

var _ = Describe("Level", func() {
	Describe("String", func() {
		When("level is explicitly unknown", func() {
			It("should return UNKNOWN", func() {
				Expect(alert.Unknown.String()).Should(Equal("UNKNOWN"))
			})
		})

		When("level is info", func() {
			It("should return INFO", func() {
				Expect(alert.Info.String()).Should(Equal("INFO"))
			})
		})

		When("level is warn", func() {
			It("should return WARN", func() {
				Expect(alert.Warn.String()).Should(Equal("WARN"))
			})
		})

		When("level is error", func() {
			It("should return ERROR", func() {
				Expect(alert.Error.String()).Should(Equal("ERROR"))
			})
		})

		When("level is unknown", func() {
			It("should return UNKNOWN", func() {
				for i := 4; i < 255; i++ {
					Expect(alert.Level(i).String()).Should(Equal("UNKNOWN"))
				}
			})
		})
	})

	Describe("MarshalJSON", func() {
		When("level is explicitly unknown", func() {
			It("should return UNKNOWN", func() {
				data, err := alert.Unknown.MarshalJSON()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).Should(Equal([]byte(`"UNKNOWN"`)))
			})
		})

		When("level is info", func() {
			It("should return INFO", func() {
				data, err := alert.Info.MarshalJSON()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).Should(Equal([]byte(`"INFO"`)))
			})
		})

		When("level is warn", func() {
			It("should return WARN", func() {
				data, err := alert.Warn.MarshalJSON()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).Should(Equal([]byte(`"WARN"`)))
			})
		})

		When("level is error", func() {
			It("should return ERROR", func() {
				data, err := alert.Error.MarshalJSON()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).Should(Equal([]byte(`"ERROR"`)))
			})
		})

		When("level is unknown", func() {
			It("should return UNKNOWN", func() {
				data, err := alert.Level(240).MarshalJSON()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).Should(Equal([]byte(`"UNKNOWN"`)))
			})
		})
	})

	Describe("UnmarshalJSON", func() {
		When("level is explicitly UNKNOWN", func() {
			It("should return Unknown", func() {
				level := alert.Level(255)
				err := level.UnmarshalJSON([]byte(`"UNKNOWN"`))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(level).Should(Equal(alert.Unknown))
			})
		})

		When("level is INFO", func() {
			It("should return Info", func() {
				level := alert.Level(255)
				err := level.UnmarshalJSON([]byte(`"INFO"`))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(level).Should(Equal(alert.Info))
			})
		})

		When("level is WARN", func() {
			It("should return Warn", func() {
				level := alert.Level(255)
				err := level.UnmarshalJSON([]byte(`"WARN"`))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(level).Should(Equal(alert.Warn))
			})
		})

		When("level is ERROR", func() {
			It("should return Error", func() {
				level := alert.Level(255)
				err := level.UnmarshalJSON([]byte(`"ERROR"`))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(level).Should(Equal(alert.Error))
			})
		})

		When("level is whatever", func() {
			It("should return Unknown", func() {
				level := alert.Level(255)
				err := level.UnmarshalJSON([]byte(`"something"`))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(level).Should(Equal(alert.Unknown))
			})
		})
	})

})
