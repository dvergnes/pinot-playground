package main

import (
	"github.com/dvergnes/pinot-playground/cert-monitor/alert"
	"github.com/dvergnes/pinot-playground/cert-monitor/internal/version"
	"log"
	"net/http"
	"time"

	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"

	"go.uber.org/zap"
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

	suggaredLogger.Info("starting certificate monitor", "version", version.Version)

	// 1. read config
	config := monitor.Config{}
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
