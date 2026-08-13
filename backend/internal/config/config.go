package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT      string
	DB_URL    string
	JWTSECRET string
	BASEURL   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found")
	}
	
	cfg := &Config{
		PORT:      getEnv("PORT", "8080"),
		DB_URL:    os.Getenv("DATABASE_URL"),
		JWTSECRET: os.Getenv("JWT_SECRET"),
		BASEURL:   getEnv("BASE_URL", "http://localhost:8080"),
	}

	if cfg.DB_URL == "" {
		return nil, fmt.Errorf("Database url is required")
	}

	if cfg.JWTSECRET == "" {
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
