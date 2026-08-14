package main

import (
	"log"
	"net/http"

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

	shortenerRepo := shortener.NewPostgresRepository(pool)
	shortenerService := shortener.NewService(shortenerRepo)

	shortenerHandler := shortener.NewHandler(shortenerService, noopClickLogger{}, cfg.BASEURL)

	router := httpserver.NewRouter(shortenerHandler)

	log.Printf("listening on :%s", cfg.PORT)

	if err := http.ListenAndServe(":"+cfg.PORT, router); err != nil {
		log.Fatal(err)
	}
}

type noopClickLogger struct{}

func (noopClickLogger) LogClick(code, referrer, userAgent string) {}
