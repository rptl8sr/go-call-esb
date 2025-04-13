package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Get sends an HTTP GET request to the specified URL with the provided headers and returns a Response or an error.
func (c *client) Get(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.do(ctx, http.MethodGet, url, headers, nil)
}

func (c *client) Post(ctx context.Context, url string, headers map[string]string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPost, url, headers, body)
}

// Put sends an HTTP PUT request to the specified URL with the given headers and body, returning a Response or an error.
func (c *client) Put(ctx context.Context, url string, headers map[string]string, body interface{}) (*Response, error) {
	return c.do(ctx, http.MethodPut, url, headers, body)
}

// Delete sends an HTTP DELETE request to the specified URL with the provided headers and returns a Response or an error.
func (c *client) Delete(ctx context.Context, url string, headers map[string]string) (*Response, error) {
	return c.do(ctx, http.MethodDelete, url, headers, nil)
}

// do execute an HTTP request with the given method, URL, headers, and body, handling retries and logging.
func (c *client) do(ctx context.Context, method, url string, headers map[string]string, body interface{}) (*Response, error) {
	var lastErr error
	var retries int

	for i := 0; i <= c.config.MaxRetries; i++ {
		startTime := time.Now()

		// Логируем запрос
		c.config.Logger.Request(ctx, method, url, headers, body, startTime)

		resp, err := c.execute(ctx, method, url, headers, body)
		duration := time.Since(startTime)

		if err != nil {
			lastErr = err
			c.config.Logger.Response(ctx, method, url, 0, nil, 0, duration, err)

			if i < c.config.MaxRetries {
				retries++
				c.config.Logger.Retry(ctx, method, url, i+1, err)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(c.config.RetryDelay):
					continue
				}
			}
			continue
		}

		// Convert http.Header to map[string]string for logger
		headerMap := make(map[string]string)
		for k, v := range resp.Headers {
			if len(v) > 0 {
				headerMap[k] = v[0]
			}
		}
		c.config.Logger.Response(ctx, method, url, resp.StatusCode, headerMap, len(resp.Body), duration, nil)

		resp.Retries = retries
		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// execute выполняет один HTTP запрос
func (c *client) execute(ctx context.Context, method, url string, headers map[string]string, body interface{}) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	for k, v := range c.config.DefaultHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
		Headers:    resp.Header,
	}, nil
}
