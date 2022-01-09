package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dvergnes/pinot-playground/cert-monitor/alert"
	"github.com/dvergnes/pinot-playground/cert-monitor/internal/version"
	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
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
	config,err := newFromCLI()
	if err!=nil {
		suggaredLogger.Fatalw("failed to initialize application", "error", err)
	}
	// 2. init app
	httpClient := &http.Client{
		Timeout: config.Scrapping.Timeout,
	}
	gatherer := monitor.NewPrometheusCertificateInfosGatherer(
		suggaredLogger.Named("certInfoGatherer"),
		httpClient,
		config.Scrapping.Endpoint,
		config.Scrapping.Metric)
	notifier := alert.NewLogNotifier(suggaredLogger.Named("logNotifier"))
	certMonitor := monitor.NewCertificateMonitor(
		suggaredLogger.Named("monitor"),
		gatherer,
		notifier,
		sysClock,
		config.Threshold)
	// 3. run the monitor
	if err := certMonitor.CheckCertificates(); err != nil {
		suggaredLogger.Fatalw("failed to verify certificate", "error", err)
	}

}

func newFromBytes(data []byte) (*monitor.Config, error) {
	config := monitor.Config{

	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %s", err)
	}

	// TODO: config validation
	return &config, nil
}

func newFromFile(filepath string) (*monitor.Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf( "failed to read config file: %s", err)
	}
	return newFromBytes(data)
}

func newFromCLI() (*monitor.Config, error) {
	configPath := flag.String("config", "", "Configuration file path")
	flag.Parse()

	if *configPath == "" {
		return nil, errors.New("missing config file path. Refer --help")
	}

	return newFromFile(*configPath)
}
