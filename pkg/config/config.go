package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Host:        getEnv("HOST", "0.0.0.0"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (cfg *Config) GetAddr() string {
	return cfg.Host + ":" + cfg.Port
}
