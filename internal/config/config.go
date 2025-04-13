package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Mode string

const (
	Prod Mode = "prod"
	Dev  Mode = "dev"
)

// Config is a composite configuration struct aggregating settings for the application.
type Config struct {
	App
	ESB
	YCF
	TG
}

// The App represents the core application configuration with runtime, logging, and concurrency settings.
type App struct {
	Mode            Mode       `env:"APP_MODE" envDefault:"dev"`
	LogLevel        slog.Level `env:"APP_LOG_LEVEL" envDefault:"info"`
	MaxRoutines     int        `env:"APP_MAX_GOROUTINES" envDefault:"5"`
	TimeoutDuration int        `env:"APP_TIMEOUT_DURATION" envDefault:"60"`
}

// ESB contains configuration for the Enterprise Service Bus.
type ESB struct {
	BaseURL   *url.URL      `env:"ESB_BASE_URL"`
	APIKey    string        `env:"ESB_API_KEY"`
	Path      string        `env:"ESB_PATH"`
	TimeDelta time.Duration `env:"ESB_TIME_DELTA"`
}

type TG struct {
	APIKey string `env:"TG_API_KEY"`
	UserID int    `env:"TG_USER_ID"`
}

// YCF holds a configuration for a Yandex Cloud Functions.
type YCF struct {
	SaID     string `env:"YCF_SA_ID"`
	Cron     string `env:"YCF_CRON"`
	FuncName string `env:"YCF_FUNC_NAME"`
}

// Must load the configuration and panics if it fails.
// Use this when configuration is required for the application to start.
func Must() Config {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		panic(fmt.Sprintf("Error processing environment variables: %v", err))
	}

	return config
}
