package analytics

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/stats/{code}", h.handleGetStats)
}

func (h *Handler) handleGetStats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	stats, err := h.service.Stats(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
