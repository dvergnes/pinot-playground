package monitor

import (
	"fmt"
	"os"
	"time"

	"github.com/dvergnes/pinot-playground/cert-monitor/alert"

	"go.uber.org/zap"
)

// CertificateInfo contains the name, namespace and the expiration of a certificate declared in k8s
type CertificateInfo struct {
	// Name of the certificate in k8s
	Name string
	// Namespace where the certificate is defined
	Namespace string
	// Expiration defines the timestamp of when the certificate will expire in nanoseconds since the epoch
	Expiration int64
}

// CertificateInfoGatherer collects information about the certificates defined in k8s
type CertificateInfoGatherer interface {
	GatherCertificateInfos() ([]CertificateInfo, error)
}

// TODO:doc
type Clock interface {
	Now() int64
}

func NewCertificateMonitor(logger *zap.SugaredLogger, gatherer CertificateInfoGatherer, notifier alert.Notifier, clock Clock, threshold time.Duration) *CertificateMonitor {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Warn("failed to determine hostname, using unknown value")
		hostname = "unknown"
	}
	return &CertificateMonitor{
		hostname:                hostname,
		clock:                   clock,
		threshold:               threshold.Nanoseconds(),
		certificateInfoGatherer: gatherer,
		notifier:                notifier,
	}
}

// TODO:doc
type CertificateMonitor struct {
	hostname  string
	threshold int64

	clock                   Clock
	certificateInfoGatherer CertificateInfoGatherer
	notifier                alert.Notifier
}

func (cm *CertificateMonitor) CheckCertificates() error {
	certInfos, err := cm.certificateInfoGatherer.GatherCertificateInfos()
	if err != nil {
		return fmt.Errorf("failed to gather certificate information: %s", err)
	}

	var alerts []alert.Alert
	for _, cert := range certInfos {
		now := cm.clock.Now()
		delta := cert.Expiration - now
		if delta <= 0 {
			alerts = append(alerts, alert.Alert{
				Level: alert.Error,
				ObjectRef: alert.ObjectRef{
					Namespace: cert.Namespace,
					Name:      cert.Name,
				},
				Message: "certificate expired",
				When:    now,
				Source:  cm.hostname,
			})
		} else if delta <= cm.threshold {
			alerts = append(alerts, alert.Alert{
				Level: alert.Warn,
				ObjectRef: alert.ObjectRef{
					Namespace: cert.Namespace,
					Name:      cert.Name,
				},
				Message: "certificate is about to expire",
				When:    now,
				Source:  cm.hostname,
			})
		}
	}

	return cm.notify(alerts)
}

func (cm *CertificateMonitor) notify(b []alert.Alert) error {
	for _, a := range b {
		if err := cm.notifier.Send(a); err != nil {
			return fmt.Errorf("failed to send a for certificate %s.%s: %s",
				a.ObjectRef.Namespace, a.ObjectRef.Name,
				err)
		}
	}
	return nil
}
