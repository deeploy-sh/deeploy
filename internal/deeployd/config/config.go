package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv        string
	Port         string
	DBDriver     string
	DBConnection string
	JWTSecret    string
	CookieSecure bool
}

func Load() *Config {
	godotenv.Load()

	goEnv := getEnv("GO_ENV", "development")

	return &Config{
		GoEnv:        goEnv,
		Port:         getEnv("PORT", "8090"),
		DBDriver:     getEnv("DB_DRIVER", "sqlite"),
		DBConnection: getEnv("DB_CONNECTION", "./data/deeploy.db"),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		CookieSecure: goEnv != "development",
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
