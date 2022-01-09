package monitor_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("PrometheusScrapper", func() {
	const (
		endpoint   = "http://localhost:8080/metrics"
		metricName = "certmanager_certificate_expiration_timestamp_seconds"
	)
	var (
		scrapper monitor.CertificateInfoGatherer
		logger   = zap.S()
	)

	Describe("GatherCertificateInfos", func() {

		var (
			certs []monitor.CertificateInfo
			err   error
		)

		JustBeforeEach(func() {
			certs, err = scrapper.GatherCertificateInfos()
		})

		When("metrics can be collected", func() {
			BeforeEach(func() {
				client := NewTestClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					Expect(req.URL.String()).Should(Equal(endpoint))
					return &http.Response{
						StatusCode: http.StatusOK,
						// Send response to be tested
						Body: ioutil.NopCloser(bytes.NewBufferString(`# HELP certmanager_certificate_expiration_timestamp_seconds The date after which the certificate expires. Expressed as a Unix Epoch Time.
# TYPE certmanager_certificate_expiration_timestamp_seconds gauge
certmanager_certificate_expiration_timestamp_seconds{name="example-com",namespace="sandbox"} 1.649444639e+09
certmanager_certificate_expiration_timestamp_seconds{name="my-selfsigned-ca",namespace="sandbox"} 1.649444074e+09
certmanager_certificate_expiration_timestamp_seconds{name="vergnes-com",namespace="sandbox"} 1.641709414e+09
# HELP certmanager_certificate_ready_status The ready status of the certificate.
# TYPE certmanager_certificate_ready_status gauge
certmanager_certificate_ready_status{condition="False",name="example-com",namespace="sandbox"} 0
certmanager_certificate_ready_status{condition="False",name="my-selfsigned-ca",namespace="sandbox"} 0
certmanager_certificate_ready_status{condition="False",name="vergnes-com",namespace="sandbox"} 0
certmanager_certificate_ready_status{condition="True",name="example-com",namespace="sandbox"} 1
certmanager_certificate_ready_status{condition="True",name="my-selfsigned-ca",namespace="sandbox"} 1
certmanager_certificate_ready_status{condition="True",name="vergnes-com",namespace="sandbox"} 1
certmanager_certificate_ready_status{condition="Unknown",name="example-com",namespace="sandbox"} 0
certmanager_certificate_ready_status{condition="Unknown",name="my-selfsigned-ca",namespace="sandbox"} 0
certmanager_certificate_ready_status{condition="Unknown",name="vergnes-com",namespace="sandbox"} 0
`)),
						// Must be set to non-nil value or it panics
						Header: map[string][]string{
							"Content-Type": {"text/plain; version=0.0.4; charset=utf-8"},
						},
					}, nil
				})
				scrapper = monitor.NewPrometheusCertificateInfosGatherer(logger, client, endpoint, metricName)
			})

			It("should read the cerificate info from the metrics", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(certs).Should(HaveLen(3))
				Expect(certs).Should(ContainElements(
					monitor.CertificateInfo{
						Name:       "example-com",
						Namespace:  "sandbox",
						Expiration: 1649444639000000000,
					},
					monitor.CertificateInfo{
						Name:       "my-selfsigned-ca",
						Namespace:  "sandbox",
						Expiration: 1649444074000000000,
					},
					monitor.CertificateInfo{
						Name:       "vergnes-com",
						Namespace:  "sandbox",
						Expiration: 1641709414000000000,
					},
				))
			})
		})

		When("endpoint is down", func() {
			var criticalError = errors.New("host unreachable")
			BeforeEach(func() {
				client := NewTestClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					Expect(req.URL.String()).Should(Equal(endpoint))
					return nil, criticalError
				})
				scrapper = monitor.NewPrometheusCertificateInfosGatherer(logger, client, endpoint, metricName)
			})

			It("should propagate the error", func() {
				Expect(err).Should(MatchError("error making HTTP request to: Get \"http://localhost:8080/metrics\": host unreachable"))
			})

		})

		When("endpoint returns a HTTP error", func() {
			BeforeEach(func() {
				client := NewTestClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					Expect(req.URL.String()).Should(Equal(endpoint))
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
					}, nil
				})
				scrapper = monitor.NewPrometheusCertificateInfosGatherer(logger, client, endpoint, metricName)
			})

			It("should propagate the error", func() {
				Expect(err).Should(MatchError("http://localhost:8080/metrics returned HTTP status 500"))
			})
		})

		When("metrics cannot be parsed", func() {
			BeforeEach(func() {
				client := NewTestClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					Expect(req.URL.String()).Should(Equal(endpoint))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString("invalid payload\n")),
					}, nil
				})
				scrapper = monitor.NewPrometheusCertificateInfosGatherer(logger, client, endpoint, metricName)
			})
			It("should propagate the error", func() {
				Expect(err).Should(MatchError("error getting processing metrics for http://localhost:8080/metrics: text format parsing error in line 1: expected float as value, got \"payload\""))
			})
		})

		When("metrics cannot be found", func() {
			BeforeEach(func() {
				client := NewTestClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					Expect(req.URL.String()).Should(Equal(endpoint))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString("\n")),
					}, nil
				})
				scrapper = monitor.NewPrometheusCertificateInfosGatherer(logger, client, endpoint, metricName)
			})
			It("should return an empty slice", func() {
				Expect(certs).Should(BeEmpty())
			})
		})

		When("metrics has an unexpected type", func() {
			BeforeEach(func() {
				client := NewTestClient(func(req *http.Request) (*http.Response, error) {
					// Test request parameters
					Expect(req.URL.String()).Should(Equal(endpoint))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(`# HELP certmanager_certificate_expiration_timestamp_seconds The date after which the certificate expires. Expressed as a Unix Epoch Time.
# TYPE certmanager_certificate_expiration_timestamp_seconds counter
certmanager_certificate_expiration_timestamp_seconds{name="example-com",namespace="sandbox"} 1.649444639e+09
`)),
					}, nil
				})
				scrapper = monitor.NewPrometheusCertificateInfosGatherer(logger, client, endpoint, metricName)
			})
			It("should return an error", func() {
				Expect(certs).Should(BeEmpty())
				Expect(err).Should(MatchError("metric family certmanager_certificate_expiration_timestamp_seconds is a COUNTER and not a GAUGE"))
			})
		})

	})
})

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
