package config

import (
	"fmt"
	"os"
)

type Config struct {
	PORT      string
	DB_URL    string
	JWTSECRET string
	BASEURL   string
}

func Load() (*Config, error) {
	cfg := &Config{
		PORT: getEnv("PORT", "8080"),
		DB_URL: os.Getenv("DB_URL"),
		JWTSECRET: os.Getenv("JWTSECRET"),
		BASEURL: getEnv("BASEURL", "http://localhost:8080"),
	}

	if cfg.DB_URL == "" {
		return nil, fmt.Errorf("Database url is required")
	}

	if cfg.JWTSECRET == "" {
		return nil, fmt.Errorf("JWT secret is required")
	}

	return cfg, nil
}


func getEnv(key, fallback string,) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}