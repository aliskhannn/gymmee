package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration parameters.
type Config struct {
	Env           string
	BotToken      string
	DatabaseURL   string
	ServerAddress string
}

// MustLoad reads the .env file and populates the Config struct.
func MustLoad() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env:           getEnv("ENV", "development"),
		BotToken:      getEnv("BOT_TOKEN", ""),
		DatabaseURL:   getEnv("GOOSE_DBSTRING", "./gymlog.db"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN environment variable is required")
	}

	return cfg, nil
}

// getEnv retrieves an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
