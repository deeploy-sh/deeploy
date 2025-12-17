package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CookieSecure     bool
	AppEnv           string
	GitHubToken      string
	ResendAPIKey     string
	ResendAudienceID string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("no .env file found, using environment variables")
	}

	appEnv := envString("APP_ENV", "development")

	AppConfig = &Config{
		AppEnv:           appEnv,
		CookieSecure:     appEnv == "production",
		GitHubToken:      envString("GITHUB_TOKEN", ""),
		ResendAPIKey:     envString("RESEND_API_KEY", ""),
		ResendAudienceID: envString("RESEND_AUDIENCE_ID", ""),
	}

	// Warn if email subscription won't work (but don't crash - it's optional)
	if AppConfig.ResendAPIKey == "" || AppConfig.ResendAudienceID == "" {
		slog.Warn("RESEND_API_KEY or RESEND_AUDIENCE_ID not set, email subscription disabled")
	}
}

func envString(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return def
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}
