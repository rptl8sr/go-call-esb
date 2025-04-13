package esb

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type call struct {
	ApiKey    string `json:"api_key"`
	Timestamp int    `json:"timestamp"`
}

// Call retrieves a list of stores from the ESB API based on the provided query parameters.
func (c *Client) Call(ctx context.Context, endpoint string, t time.Time) error {
	u, err := url.JoinPath(c.baseURL.String(), endpoint)
	if err != nil {
		return fmt.Errorf("failed to join path: %w", err)
	}

	reqBody := call{
		ApiKey:    c.apiKey,
		Timestamp: int(t.Unix()),
	}

	resp, err := c.httpClient.Put(ctx, u, getHeaders(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to get stores: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
