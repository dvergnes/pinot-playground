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
