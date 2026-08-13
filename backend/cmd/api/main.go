package main

import (
	"log"

	"github.com/alex-njaiya/url-shortener/internal/config"
	"github.com/alex-njaiya/url-shortener/internal/database"
)


func main() {
	// load the configs
	// create a db pool

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	pool, err := database.NewPool(cfg.DB_URL)

	if err != nil {
		log.Fatalf("connection to database: %v", err)
	}

	defer pool.Close()

	log.Printf("listening on :%s", cfg.PORT)
}