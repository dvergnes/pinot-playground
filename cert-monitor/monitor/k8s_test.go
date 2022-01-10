package monitor_test

import (
	"context"
	"errors"
	v1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	"github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"github.com/dvergnes/pinot-playground/cert-monitor/mocks"
	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("k8sCertificateInfoGatherer", func() {

	var (
		clientSetMock   *mocks.Interface
		certManagerMock *mocks.CertmanagerV1Interface
		certAPIMock     *mocks.CertificateInterface
		gatherer        monitor.CertificateInfoGatherer
	)

	BeforeEach(func() {
		clientSetMock = &mocks.Interface{}
		certManagerMock = &mocks.CertmanagerV1Interface{}
		clientSetMock.On("CertmanagerV1").Return(certManagerMock)
		certAPIMock = &mocks.CertificateInterface{}
		certManagerMock.On("Certificates", "").Return(certAPIMock)
		gatherer = monitor.NewKubernetesCertificateInfoGatherer(zap.S(), clientSetMock, monitor.GathererConfig{
			PageSize: 1,
			Timeout:  time.Second,
		})
	})

	AfterEach(func() {
		clientSetMock.AssertExpectations(GinkgoT())
		certManagerMock.AssertExpectations(GinkgoT())
		certAPIMock.AssertExpectations(GinkgoT())
	})

	Describe("GatherCertificateInfos", func() {

		var (
			certs []monitor.CertificateInfo
			err   error
		)

		JustBeforeEach(func() {
			certs, err = gatherer.GatherCertificateInfos(context.TODO())
		})
		When("no pagination", func() {
			expiry := metav1.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
			BeforeEach(func() {
				certList := &v1.CertificateList{
					Items: []v1.Certificate{
						{
							Status: v1.CertificateStatus{
								NotAfter: &expiry,
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "cert",
								Namespace: "ns",
							},
						},
					},
				}
				certAPIMock.On("List",
					mock.AnythingOfType("*context.timerCtx"),
					mock.MatchedBy(func(opts metav1.ListOptions) bool {
						Expect(opts.Limit).Should(BeEquivalentTo(1))
						return Expect(opts.Continue).Should(BeEmpty())
					})).Once().Return(certList, nil)
			})
			It("should return all certificates", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(certs).Should(HaveLen(1))
				Expect(certs).Should(ContainElements(monitor.CertificateInfo{
					Namespace:  "ns",
					Name:       "cert",
					Expiration: expiry.UnixNano(),
				}))
			})

		})

		When("pagination", func() {
			expiry := metav1.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
			BeforeEach(func() {
				page1 := &v1.CertificateList{
					Items: []v1.Certificate{
						{
							Status: v1.CertificateStatus{
								NotAfter: &expiry,
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "cert1",
								Namespace: "ns1",
							},
						},
					},
					ListMeta: metav1.ListMeta{
						Continue: "go-to-page-2",
					},
				}
				page2 := &v1.CertificateList{
					Items: []v1.Certificate{
						{
							Status: v1.CertificateStatus{
								NotAfter: &expiry,
							},
							ObjectMeta: metav1.ObjectMeta{
								Name:      "cert2",
								Namespace: "ns2",
							},
						},
					},
				}
				pageID := 0
				certAPIMock.On("List",
					mock.AnythingOfType("*context.timerCtx"),
					mock.MatchedBy(func(opts metav1.ListOptions) bool {
						Expect(opts.Limit).Should(BeEquivalentTo(1))
						return Expect(opts.Continue).Should(Or(BeEmpty(),Equal("go-to-page-2")))
					})).
					Twice().
					Return(func(context.Context, metav1.ListOptions) *v1.CertificateList {
						if pageID == 0 {
							pageID++
							return page1
						}
						return page2
					}, nil)
			})
			It("should return all certificates", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(certs).Should(HaveLen(2))
				Expect(certs).Should(ContainElements(
					monitor.CertificateInfo{
						Namespace:  "ns1",
						Name:       "cert1",
						Expiration: expiry.UnixNano(),
					},
					monitor.CertificateInfo{
						Name:       "cert2",
						Namespace:  "ns2",
						Expiration: expiry.UnixNano(),
					}))
			})
		})

		When("API call failed", func() {
			criticalError := errors.New("failed to call API")
			BeforeEach(func() {
				certAPIMock.On("List", mock.Anything, mock.Anything).Return(nil, criticalError)
			})
			It("should propagate the error", func() {
				Expect(err).Should(MatchError("failed to fetch certificates: failed to call API"))
			})
		})
	})

})
