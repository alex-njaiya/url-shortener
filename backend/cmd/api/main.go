package main

import (
	"log"
	"net/http"

	"github.com/alex-njaiya/url-shortener/internal/analytics"
	"github.com/alex-njaiya/url-shortener/internal/auth"
	"github.com/alex-njaiya/url-shortener/internal/config"
	"github.com/alex-njaiya/url-shortener/internal/database"
	"github.com/alex-njaiya/url-shortener/internal/httpserver"
	"github.com/alex-njaiya/url-shortener/internal/shortener"
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

	// ANALYTICS
	analyticsRepo := analytics.NewPostgresRepository(pool)
	analyticsService := analytics.NewService(analyticsRepo)

	analyticsHandler := analytics.NewHandler(analyticsService)

	// URL SHORTENER
	shortenerRepo := shortener.NewPostgresRepository(pool)
	shortenerService := shortener.NewService(shortenerRepo)

	shortenerHandler := shortener.NewHandler(shortenerService, analyticsService, cfg.BASEURL, cfg.JWT_SECRET)

	// USER REGISTER AND LOGIN
	authRepo := auth.NewPostgresRepository(pool)
	authService := auth.NewService(authRepo)

	authHandler := auth.NewHandler(authService, cfg.JWT_SECRET)

	router := httpserver.NewRouter(shortenerHandler, authHandler, analyticsHandler)

	log.Printf("listening on :%s", cfg.PORT)

	if err := http.ListenAndServe(":"+cfg.PORT, router); err != nil {
		log.Fatal(err)
	}
}
