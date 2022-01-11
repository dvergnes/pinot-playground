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

package monitor

import (
	"context"
	"fmt"

	certmanager "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"go.uber.org/zap"
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
