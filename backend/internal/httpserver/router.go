package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RouteRegistrar interface {
	RegisterRoutes(r chi.Router)
}

func NewRouter(registrars ...RouteRegistrar) http.Handler {
	r := chi.NewRouter()

	// middlewares
	r.Use(CORS)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	for _, reg := range registrars {
		reg.RegisterRoutes(r)
	}

	return r
}
