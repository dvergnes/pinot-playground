package monitor

import "time"

// Config contains the configuration for the monitor
type Config struct {
	// Threshold defines the duration for which a certificate is not considered as close to expiration. For example,
	// if a certificate is valid for the next 20 days but the threshold is set to 30 days, it is considered as close to
	// expiration.
	Threshold time.Duration `yaml:"threshold"`
	// GathererConfig contains the configuration for fetching the certificate info
	GathererConfig GathererConfig `yaml:"gatherer"`
}

// GathererConfig contains the configuration for fetching the certificate info
type GathererConfig struct {
	// PageSize defines the page size when calling the list certificate API
	PageSize int64          `yaml:"page_size"`
	// Timeout defines the timeout to fetch a page
	Timeout  time.Duration `yaml:"timeout"`
}
