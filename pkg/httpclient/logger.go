package httpclient

import (
	"context"
	"log/slog"
	"time"
)

// SlogLogger реализация логгера с использованием slog
type SlogLogger struct {
	logger *slog.Logger
}

// NewLogger создает новый экземпляр SlogLogger
func NewLogger(logger *slog.Logger) *SlogLogger {
	if logger == nil {
		logger = slog.Default()
	}
	return &SlogLogger{logger: logger}
}

// Request логирует информацию о запросе
func (l *SlogLogger) Request(ctx context.Context, method, url string, headers map[string]string, body interface{}, startTime time.Time) {
	l.logger.InfoContext(ctx, "http request",
		"method", method,
		"url", url,
		"headers", headers,
		"body", body,
		"time", startTime,
	)
}

// Response логирует информацию об ответе
func (l *SlogLogger) Response(ctx context.Context, method, url string, statusCode int, headers map[string]string, bodyLength int, duration time.Duration, err error) {
	attrs := []any{
		"method", method,
		"url", url,
		"status_code", statusCode,
		"headers", headers,
		"duration", duration,
	}

	if err != nil {
		attrs = append(attrs, "error", err)
		l.logger.ErrorContext(ctx, "http response error", attrs...)
	} else {
		attrs = append(attrs, "body_length", bodyLength)
		l.logger.InfoContext(ctx, "http response", attrs...)
	}
}

// Retry логирует информацию о повторной попытке
func (l *SlogLogger) Retry(ctx context.Context, method, url string, attempt int, err error) {
	l.logger.WarnContext(ctx, "http retry",
		"method", method,
		"url", url,
		"attempt", attempt,
		"error", err,
	)
}
