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
	"github.com/dvergnes/pinot-playground/cert-monitor/internal/version"
	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"

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
		config.GathererConfig)
	notifier := alert.NewLogNotifier(suggaredLogger.Named("logNotifier"))
	certMonitor := monitor.NewCertificateMonitor(
		suggaredLogger.Named("monitor"),
		gatherer,
		notifier,
		sysClock,
		config.Threshold)
	// 3. run the monitor
	if err := certMonitor.CheckCertificates(context.Background()); err != nil {
		suggaredLogger.Fatalw("failed to verify certificate", "error", err)
	}

}

func newFromBytes(data []byte) (*monitor.Config, error) {
	config := monitor.Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %s", err)
	}

	// TODO: config validation
	return &config, nil
}

func newFromFile(filepath string) (*monitor.Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}
	return newFromBytes(data)
}

func newFromCLI(logger *zap.SugaredLogger) (*monitor.Config, *rest.Config, error) {
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
