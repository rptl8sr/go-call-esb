package httpclient

import (
	"time"
)

// Config configuration of an HTTP client
type Config struct {
	Timeout            time.Duration
	MaxRetries         int
	RetryDelay         time.Duration
	BaseURL            string
	DefaultHeaders     map[string]string
	InsecureSkipVerify bool
	Logger             *SlogLogger
}

// DefaultConfig return default config
func DefaultConfig() *Config {
	return &Config{
		Timeout:            30 * time.Second,
		MaxRetries:         3,
		RetryDelay:         1 * time.Second,
		DefaultHeaders:     make(map[string]string),
		InsecureSkipVerify: false,
		Logger:             NewLogger(nil),
	}
}
