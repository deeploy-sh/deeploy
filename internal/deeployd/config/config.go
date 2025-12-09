package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv         string
	Port          string
	DatabaseURL   string
	JWTSecret     string
	EncryptionKey string
	CookieSecure  bool
}

func Load() *Config {
	godotenv.Load()

	goEnv := getEnv("GO_ENV", "development")
	isProduction := goEnv == "production"

	databaseURL := requireEnv("DATABASE_URL", "postgres://user:pass@host:5432/deeploy", true)
	jwtSecret := requireEnv("JWT_SECRET", "min 32 characters", isProduction)
	encryptionKey := requireEnv("ENCRYPTION_KEY", "exactly 32 characters", isProduction)

	return &Config{
		GoEnv:         goEnv,
		Port:          getEnv("PORT", "8090"),
		DatabaseURL:   databaseURL,
		JWTSecret:     jwtSecret,
		EncryptionKey: encryptionKey,
		CookieSecure:  isProduction,
	}
}

func (c *Config) IsDevelopment() bool {
	return c.GoEnv == "development"
}

func (c *Config) IsProduction() bool {
	return c.GoEnv == "production"
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func requireEnv(key, hint string, required bool) string {
	value := os.Getenv(key)
	if value == "" && required {
		fmt.Printf("\n Missing required environment variable: %s\n", key)
		fmt.Printf("   Expected: %s\n\n", hint)
		os.Exit(1)
	}
	return value
}
