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

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dvergnes/pinot-playground/cert-monitor/alert"
	"github.com/dvergnes/pinot-playground/cert-monitor/config"
	"github.com/dvergnes/pinot-playground/cert-monitor/internal/version"
	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"

	"github.com/Shopify/sarama"
	certmanager "github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var sysClock = &systemClock{}

type systemClock struct {
}

func (sc *systemClock) Now() int64 {
	return time.Now().UnixNano()
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger %+v", err)
	}
	defer logger.Sync() // flushes buffer, if any

	suggaredLogger := logger.Sugar()

	suggaredLogger.Infow("starting certificate monitor", "version", version.Version)

	// 1. read config
	config, k8sCfg, err := newFromCLI(suggaredLogger)
	if err != nil {
		suggaredLogger.Fatalw("failed to initialize application", "error", err)
	}
	// 2. init app
	clientSet, err := certmanager.NewForConfig(k8sCfg)
	if err != nil {
		suggaredLogger.Fatalw("failed to create k8s client", "error", err)
	}
	gatherer := monitor.NewKubernetesCertificateInfoGatherer(
		suggaredLogger.Named("k8sCertInfoGatherer"),
		clientSet,
		config.Monitor.GathererConfig)
	//notifier := alert.NewLogNotifier(suggaredLogger.Named("logNotifier"))
	producer,err := sarama.NewSyncProducer(config.Notifier.Brokers, nil)
	if err!=nil {
		suggaredLogger.Fatalw("failed to create kafka producer", "error", err)
	}
	notifier := alert.NewKafkaNotifier(config.Notifier.Topic, producer)
	certMonitor := monitor.NewCertificateMonitor(
		suggaredLogger.Named("monitor"),
		gatherer,
		notifier,
		sysClock,
		config.Monitor.Threshold)
	// 3. run the monitor
	if err := certMonitor.CheckCertificates(context.Background()); err != nil {
		suggaredLogger.Fatalw("failed to verify certificate", "error", err)
	}

}

func newFromBytes(data []byte) (*config.Config, error) {
	config := config.Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %s", err)
	}

	// TODO: config validation
	return &config, nil
}

func newFromFile(filepath string) (*config.Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}
	return newFromBytes(data)
}

func newFromCLI(logger *zap.SugaredLogger) (*config.Config, *rest.Config, error) {
	configPath := flag.String("config", "", "Configuration file path")
	kubeConfigPath := flag.String("kubeconfig", "", "Kubectl configuration file path")
	flag.Parse()

	if *configPath == "" {
		return nil, nil, errors.New("missing config file path. Refer --help")
	}
	cfg, err := newFromFile(*configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	var k8sCfg *rest.Config
	if *kubeConfigPath == "" {
		logger.Info("creating k8s client using in-cluster configuration")
		k8sCfg, err = rest.InClusterConfig()
	} else {
		logger.Infow("creating k8s client using external configuration", "path", *kubeConfigPath)
		k8sCfg, err = clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return cfg, k8sCfg, nil
}
