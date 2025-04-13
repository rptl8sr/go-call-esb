package esb

import (
	"log/slog"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name            string
		baseURL         string
		apiKey          string
		timeout         *time.Duration
		expectedTimeout time.Duration
		logger          *slog.Logger
	}{
		{
			name:            "Valid Configuration",
			baseURL:         "https://example.com/api/",
			apiKey:          "testapikey",
			timeout:         func() *time.Duration { d := 5 * time.Second; return &d }(),
			expectedTimeout: 5 * time.Second,
			logger:          slog.Default(),
		},
		{
			name:            "Nil Timeout",
			baseURL:         "https://example.com/api/",
			apiKey:          "testapikey",
			timeout:         nil,
			expectedTimeout: defaultTimeout,
			logger:          slog.Default(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedBaseURL, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("Failed to parse baseURL: %v", err)
			}

			client := New(parsedBaseURL, tt.apiKey, tt.timeout, tt.logger)

			if client.baseURL.String() != tt.baseURL {
				t.Errorf("NewClient().baseURL = %v, expected %v", client.baseURL, tt.baseURL)
			}

			if client.apiKey != tt.apiKey {
				t.Errorf("NewClient().apiKey = %v, expected %v", client.apiKey, tt.apiKey)
			}

			if client.logger != tt.logger {
				t.Errorf("NewClient().logger = %v, expected %v", client.logger, tt.logger)
			}
		})
	}
}
