package config

import (
	"github.com/dvergnes/pinot-playground/cert-monitor/alert"
	"github.com/dvergnes/pinot-playground/cert-monitor/monitor"
)

type Config struct {
	Monitor monitor.Config `yaml:"monitor"`
	Notifier alert.KafkaConfig `yaml:"notifier"`
}
