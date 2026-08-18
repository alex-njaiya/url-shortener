package shortener

import (
	"encoding/json"
	"net/http"

	"github.com/alex-njaiya/url-shortener/internal/auth"
	"github.com/go-chi/chi/v5"
)

type ClickLogger interface {
	LogClick(code, referrer, userAgent string)
}

type Handler struct {
	service *Service
	clicks  ClickLogger
	baseURL string
}

func NewHandler(service *Service, clicks ClickLogger, baseURL string) *Handler {
	return &Handler{
		service: service,
		clicks:  clicks,
		baseURL: baseURL,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/api/shorten", h.handleShorten)
	r.Get("/{code}", h.handleRedirect)
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	Original  string `json:"original_url"`
}

func (h *Handler) handleShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	userid, exists := auth.UserIDFromContext(r.Context())

	if !exists {
		http.Error(w, "userId not found", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Example guard clause
	if h.service == nil {
		http.Error(w, "Internal server error: missing dependency", http.StatusInternalServerError)
		return
	}

	// u, err := h.service.ShortenUsingBase62Encode(r.Context(), req.URL)
	u, err := h.service.ShortenByHashing(r.Context(), &userid, req.URL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &shortenResponse{
		ShortCode: u.ShortCode,
		ShortURL:  h.baseURL + "/" + u.ShortCode,
		Original:  u.OriginalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	u, err := h.service.Resolve(r.Context(), code)

	if err != nil {
		http.Error(w, "short url not found", http.StatusNotFound)
		return
	}

	go h.clicks.LogClick(code, r.Referer(), r.UserAgent())
	http.Redirect(w, r, u.OriginalURL, http.StatusFound)
}
