package monitor

import "time"

type Config struct {
	Threshold time.Duration `yaml:"threshold"`
	GathererConfig GathererConfig `yaml:"gatherer"`
}

type GathererConfig struct {
	PageSize int64          `yaml:"page_size"`
	Timeout  time.Duration `yaml:"timeout"`
}
