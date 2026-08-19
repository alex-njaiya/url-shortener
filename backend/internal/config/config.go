package config

import (
	"fmt"
	"os"
)

type Config struct {
	PORT        string
	DB_URL      string
	JWT_SECRET  string
	BASEURL     string
	FRONTENDURL string
	IS_PROD     bool
}

func Load() (*Config, error) {
	cfg := &Config{
		PORT:        getEnv("PORT", "8080"),
		DB_URL:      os.Getenv("DATABASE_URL"),
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
		BASEURL:     getEnv("BASE_URL", "http://localhost:8080"),
		FRONTENDURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		IS_PROD:     getEnv("APP_ENV", "development") == "production",
	}

	if cfg.DB_URL == "" {
		return nil, fmt.Errorf("Database url is required")
	}
	if cfg.JWT_SECRET == "" {
		return nil, fmt.Errorf("JWT secret is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}