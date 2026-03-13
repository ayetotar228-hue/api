package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        int
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:        getEnv("HOST", "0.0.0.0"),
		Port:        port,
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
