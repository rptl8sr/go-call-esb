package telegram

import (
	"log/slog"
	"time"

	"go-call-esb/pkg/httpclient"
)

var (
	baseURL        = "https://api.telegram.org/bot"
	defaultTimeout = 10 * time.Second
)

type Client struct {
	token      string
	chatID     int
	httpClient httpclient.Client
	logger     *slog.Logger
}

func New(token string, chatID int, timeout *time.Duration, logger *slog.Logger) *Client {
	if timeout != nil {
		defaultTimeout = *timeout
	}

	return &Client{
		token:  token,
		chatID: chatID,
		httpClient: httpclient.New(&httpclient.Config{
			Timeout: defaultTimeout,
		}),
		logger: logger,
	}
}
