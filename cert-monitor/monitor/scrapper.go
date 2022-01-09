package monitor

import (
	"fmt"
	"io"
	"math"
	"net/http"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// NewPrometheusCertificateInfosGatherer creates a CertificateInfoGatherer by scrapping Prometheus metrics
func NewPrometheusCertificateInfosGatherer(client *http.Client, endpoint, metricName string) CertificateInfoGatherer {
	return &PrometheusScrapper{
		endpoint:   endpoint,
		metricName: metricName,
		httpClient: client,
	}
}

// PrometheusScrapper implements CertificateInfoGatherer by scrapping the cert-manager metrics prometheus endpoint
type PrometheusScrapper struct {
	endpoint   string
	metricName string

	httpClient *http.Client
}

// GatherCertificateInfos implements CertificateInfoGatherer contract
func (ps *PrometheusScrapper) GatherCertificateInfos() ([]CertificateInfo, error) {
	metrics, err := ps.scrapCertificateMetrics()
	if err != nil {
		return nil, err
	}
	var certs []CertificateInfo
	for _, metric := range metrics.GetMetric() {
		// TODO, check on type
		name, ns := extractNameAndNamespace(metric.GetLabel())
		certs = append(certs, CertificateInfo{
			Name:       name,
			Namespace:  ns,
			Expiration: int64(math.Round(metric.GetGauge().GetValue()))*1e9,
		})
	}
	return certs, nil
}

func (ps *PrometheusScrapper) scrapCertificateMetrics() (*dto.MetricFamily, error) {
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
