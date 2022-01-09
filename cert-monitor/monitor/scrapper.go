package monitor

import (
	"fmt"
	"io"
	"math"
	"net/http"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"go.uber.org/zap"
)

// NewPrometheusCertificateInfosGatherer creates a CertificateInfoGatherer by scrapping Prometheus metrics
func NewPrometheusCertificateInfosGatherer(logger *zap.SugaredLogger, client *http.Client, endpoint, metricName string) CertificateInfoGatherer {
	return &PrometheusScrapper{
		endpoint:   endpoint,
		metricName: metricName,
		httpClient: client,
		logger:     logger,
	}
}

// PrometheusScrapper implements CertificateInfoGatherer by scrapping the cert-manager metrics prometheus endpoint
type PrometheusScrapper struct {
	endpoint   string
	metricName string

	httpClient *http.Client
	logger     *zap.SugaredLogger
}

// GatherCertificateInfos implements CertificateInfoGatherer contract
func (ps *PrometheusScrapper) GatherCertificateInfos() ([]CertificateInfo, error) {
	metrics, err := ps.scrapCertificateMetrics()
	if err != nil {
		return nil, err
	}
	if metrics == nil {
		return nil, nil
	}
	if metrics.GetType() != dto.MetricType_GAUGE {
		return nil, fmt.Errorf("metric family %s is a %s and not a GAUGE", ps.metricName, metrics.GetType().String())
	}
	var certs []CertificateInfo
	for _, metric := range metrics.GetMetric() {
		name, ns := extractNameAndNamespace(metric.GetLabel())
		expiration := int64(math.Round(metric.GetGauge().GetValue())) * 1e9
		if name == "" || ns == "" {
			ps.logger.Warn("found certificate with empty name or namespace",
				"name", name,
				"namespace", ns,
				"expiration", expiration)
		}
		certs = append(certs, CertificateInfo{
			Name:       name,
			Namespace:  ns,
			Expiration: expiration,
		})
	}
	return certs, nil
}

func (ps *PrometheusScrapper) scrapCertificateMetrics() (*dto.MetricFamily, error) {
	ps.logger.Debug("calling metrics endpoint", "endpoint", ps.endpoint)
	resp, err := ps.httpClient.Get(ps.endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned HTTP status %d", ps.endpoint, resp.StatusCode)
	}
	format := expfmt.ResponseFormat(resp.Header)

	decoder := expfmt.NewDecoder(resp.Body, format)

	var certificateInfoMetricFamily *dto.MetricFamily

	for {
		metric := dto.MetricFamily{}
		err := decoder.Decode(&metric)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error getting processing metrics for %s: %s",
				ps.endpoint, err)
		}
		ps.logger.Debug("reading metric family", "name", metric.GetName())
		if metric.GetName() == ps.metricName {
			certificateInfoMetricFamily = &metric
		}
	}

	return certificateInfoMetricFamily, nil
}

func extractNameAndNamespace(pairs []*dto.LabelPair) (name string, ns string) {
	for _, pair := range pairs {
		if pair.GetName() == "name" {
			name = pair.GetValue()
		}
		if pair.GetName() == "namespace" {
			ns = pair.GetValue()
		}
	}
	return
}
