package monitor

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	certmanager "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewKubernetesCertificateInfoGatherer(logger *zap.SugaredLogger, clientSet certmanager.Interface, cfg GathererConfig) CertificateInfoGatherer {
	return &k8sCertificateInfoGatherer{
		cfg:       cfg,
		clientSet: clientSet,
		logger:    logger,
	}
}

type k8sCertificateInfoGatherer struct {
	cfg       GathererConfig
	clientSet certmanager.Interface

	logger *zap.SugaredLogger
}

func (k *k8sCertificateInfoGatherer) GatherCertificateInfos(parentCtx context.Context) ([]CertificateInfo, error) {
	k.logger.Info("listing certificate CRD")
	var (
		continueToken string
		page          = 1
		certInfos     []CertificateInfo
	)

	for {
		ctx, cancel := context.WithTimeout(parentCtx, k.cfg.Timeout)
		certs, err := k.clientSet.CertmanagerV1().Certificates("").List(ctx, v1.ListOptions{
			Limit:    k.cfg.PageSize,
			Continue: continueToken,
		})
		cancel()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch certificates: %w", err)
		}
		k.logger.Infow("fetched certificate CRD", "size", len(certs.Items), "page", page)
		for _, cert := range certs.Items {
			certInfos = append(certInfos, CertificateInfo{
				Name:       cert.Name,
				Namespace:  cert.Namespace,
				Expiration: cert.Status.NotAfter.UnixNano(),
			})
		}
		if certs.GetContinue() == "" {
			break
		}
		continueToken = certs.GetContinue()
		page++
	}

	return certInfos, nil
}
