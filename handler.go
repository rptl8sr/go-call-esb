package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go-call-esb/internal/clients/esb"
	"go-call-esb/internal/clients/telegram"
	"go-call-esb/internal/config"
	"go-call-esb/internal/handler"
	"go-call-esb/internal/logger"
)

// Response defines the response format for the Yandex Cloud Function.
// Used for HTTP triggers; ignored for timer triggers.
type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

// Handler is the entry point for the Yandex Cloud Function.
// Processes events from timer or HTTP triggers, fetches stores data,
// from sources, process and filters it, and create Yandex Tracker issues.
func Handler(ctx context.Context, event interface{}) (*Response, error) {
	triggerType := handler.DetectTriggerType(event)
	fmt.Printf("Trigger type: %s\n", triggerType)

	cfg := config.Must()
	log := logger.New(cfg.App.LogLevel)

	start := time.Now()
	defer func() { logger.Info("main.Handler", "time", time.Since(start).String()) }()

	if cfg.App.Mode == config.Dev && cfg.LogLevel == slog.LevelDebug {
		fmt.Println("RUNNING IN DEVELOPMENT MODE")
		fmt.Printf("config: %+v\n", cfg)
	}

	esbCLient := esb.New(
		cfg.ESB.BaseURL,
		cfg.ESB.APIKey,
		nil,
		log,
	)

	tgClient := telegram.New(
		cfg.TG.APIKey,
		cfg.TG.UserID,
		nil,
		log,
	)

	err := esbCLient.Call(
		ctx,
		cfg.ESB.Path,
		time.Now().Add(cfg.ESB.TimeDelta),
	)

	if err != nil {
		errTg := tgClient.SendMessage(
			ctx,
			fmt.Sprintf("Error calling ESB: %s", err.Error()),
		)

		if errTg != nil {
			return nil, errTg
		}
	}

	return &Response{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
