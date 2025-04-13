package esb

import (
	"log/slog"
	"net/url"
	"time"

	"go-call-esb/pkg/httpclient"
)

var (
	defaultTimeout = 30 * time.Second
)

// Client represents the main API client for interacting with the ESB system.
// It encapsulates configuration details and provides methods for performing API requests.
type Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient httpclient.Client
	logger     *slog.Logger
}

// New initializes and returns a new Client instance with the provided configuration or default settings if nil.
func New(baseURL *url.URL, apiKey string, timeout *time.Duration, logger *slog.Logger) *Client {
	if timeout != nil {
		defaultTimeout = *timeout
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: httpclient.New(&httpclient.Config{
			Timeout: defaultTimeout,
		}),
		logger: logger,
	}
}
