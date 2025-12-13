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

	// Production: validate required services
	if AppConfig.IsProduction() {
		if AppConfig.ResendAPIKey == "" {
			slog.Error("production requires RESEND_API_KEY")
			os.Exit(1)
		}
		if AppConfig.ResendAudienceID == "" {
			slog.Error("production requires RESEND_AUDIENCE_ID")
			os.Exit(1)
		}
	}
}

func envString(key, def string) string {
	if v := os.Getenv(key); v != "" {
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
