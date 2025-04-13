package telegram

import (
	"context"
	"encoding/json"
	"fmt"

	"go-report/pkg/httpclient"
)

func (c *Client) doPost(ctx context.Context, path string, body any) (*httpclient.Response, error) {
	url := fmt.Sprintf("%s%s%s", baseURL, c.token, path)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := c.httpClient.Post(ctx, url, map[string]string{
		"Content-Type": "application/json",
	}, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("failed to perform POST request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("unexpected status code from Telegram API: %d", resp.StatusCode)
	}

	return resp, nil
}
