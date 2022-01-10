package monitor_test

import (
	"errors"
	"github.com/dvergnes/pinot-playground/cert-monitor/alert"
	"github.com/stretchr/testify/mock"
	"time"

	"github.com/dvergnes/pinot-playground/cert-monitor/mocks"
	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Monitor", func() {

	const threshold = time.Minute
	var (
		gathererMock *mocks.CertificateInfoGatherer
		clockMock    *mocks.Clock
		notifierMock *mocks.Notifier
		m            *monitor.CertificateMonitor
	)

	BeforeEach(func() {
		gathererMock = &mocks.CertificateInfoGatherer{}
		clockMock = &mocks.Clock{}
		notifierMock = &mocks.Notifier{}
		m = monitor.NewCertificateMonitor(zap.S(), gathererMock, notifierMock, clockMock, threshold)
	})

	AfterEach(func() {
		gathererMock.AssertExpectations(GinkgoT())
		clockMock.AssertExpectations(GinkgoT())
		notifierMock.AssertExpectations(GinkgoT())
	})

	Describe("CheckCertificates", func() {
		var err error
		JustBeforeEach(func() {
			err = m.CheckCertificates(nil)
		})
		When("no certificates defined in the system", func() {
			BeforeEach(func() {
				gathererMock.On("GatherCertificateInfos").Return(nil, nil)
			})
			It("should not alert", func() {
				Expect(err).ShouldNot(HaveOccurred())
				// other assertions are made on the notifier mock
			})
		})

		When("certificates are valid and not close to expiration", func() {
			BeforeEach(func() {
				gathererMock.On("GatherCertificateInfos").Return([]monitor.CertificateInfo{
					{
						Name:       "cert-name",
						Namespace:  "ns",
						Expiration: time.Hour.Nanoseconds(),
					},
				}, nil)
				clockMock.On("Now").Return(int64(100))
			})
			It("should not alert", func() {
				Expect(err).ShouldNot(HaveOccurred())
				// other assertions are made on the notifier mock
			})
		})

		When("certificate is expired", func() {
			BeforeEach(func() {
				gathererMock.On("GatherCertificateInfos").Return([]monitor.CertificateInfo{
					{
						Name:       "cert-name",
						Namespace:  "ns",
						Expiration: 0,
					},
				}, nil)
				clockMock.On("Now").Return(int64(100))
				notifierMock.On("Send", mock.MatchedBy(func(a alert.Alert) bool {
					Expect(a.ObjectRef.Name).Should(Equal("cert-name"))
					Expect(a.ObjectRef.Namespace).Should(Equal("ns"))
					return Expect(a.Level).Should(Equal(alert.Error))
				})).Return(nil).Once()
			})
			It("should alert at error level", func() {
				Expect(err).ShouldNot(HaveOccurred())
				// other assertions are made on the notifier mock
			})
		})

		When("certificate is close to expiration", func() {
			BeforeEach(func() {
				now := int64(100)
				gathererMock.On("GatherCertificateInfos").Return([]monitor.CertificateInfo{
					{
						Name:       "cert-name",
						Namespace:  "ns",
						Expiration: int64(threshold) + now,
					},
				}, nil)
				clockMock.On("Now").Return(now)
				notifierMock.On("Send", mock.MatchedBy(func(a alert.Alert) bool {
					Expect(a.ObjectRef.Name).Should(Equal("cert-name"))
					Expect(a.ObjectRef.Namespace).Should(Equal("ns"))
					return Expect(a.Level).Should(Equal(alert.Warn))
				})).Return(nil).Once()
			})
			It("should alert at warn level", func() {
				Expect(err).ShouldNot(HaveOccurred())
				// other assertions are made on the notifier mock
			})
		})

		When("failed to gather certificate info", func() {
			var criticalErr = errors.New("endpoint is unreachable")
			BeforeEach(func() {
				gathererMock.On("GatherCertificateInfos").Return(nil, criticalErr)
			})
			It("should propagate the error", func() {
				Expect(err).Should(MatchError("failed to gather certificate information: endpoint is unreachable"))
			})
		})

		When("failed to send alerts", func() {
			var criticalErr = errors.New("failed to connect to SMTP server")
			BeforeEach(func() {
				gathererMock.On("GatherCertificateInfos").Return([]monitor.CertificateInfo{
					{
						Name:       "cert-name",
						Namespace:  "ns",
						Expiration: 0,
					},
				}, nil)
				clockMock.On("Now").Return(int64(100))
				notifierMock.On("Send", mock.Anything).Return(criticalErr)
			})
			It("should propagate the error", func() {
				Expect(err).Should(MatchError("failed to send an alert for certificate ns.cert-name: failed to connect to SMTP server"))
			})
		})

	})
})
