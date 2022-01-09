package monitor

import "time"

type Scrapping struct {
	Endpoint string `yaml:"endpoint"`
	Metric   string `yaml:"metric"`
	Timeout time.Duration `yaml:"timeout"`
}

type Config struct {
	Threshold time.Duration `yaml:"threshold"`
	Scrapping Scrapping `yaml:"scrapping"`
}
