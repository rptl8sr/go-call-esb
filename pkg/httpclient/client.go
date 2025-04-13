package httpclient

import (
	"context"
	"net/http"
)

// client represents an HTTP client wrapper with configurable settings and retry logic.
type client struct {
	httpClient *http.Client
	config     *Config
}

// Client defines an HTTP client interface for performing basic CRUD HTTP operations.
type Client interface {
	Get(ctx context.Context, url string, headers map[string]string) (*Response, error)
	Post(ctx context.Context, url string, headers map[string]string, body any) (*Response, error)
	Put(ctx context.Context, url string, headers map[string]string, body any) (*Response, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*Response, error)
}

// Response represents an HTTP response, including status code, body, headers, and retry count.
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Retries    int
}

// New creates a new Client instance with the provided configuration or default settings if config is nil.
func New(config *Config) Client {
	c := DefaultConfig()

	if config == nil {
		config = c
	}

	httpClient := &http.Client{
		Timeout: c.Timeout,
	}

	return &client{
		httpClient: httpClient,
		config:     c,
	}
}
