package shortener

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alex-njaiya/url-shortener/internal/auth"
	"github.com/go-chi/chi/v5"
)

type ClickLogger interface {
	LogClick(code, referrer, userAgent string)
}

type Handler struct {
	service   *Service
	clicks    ClickLogger
	baseURL   string
	jwtSecret string
}

func NewHandler(service *Service, clicks ClickLogger, baseURL, jwtSecret string) *Handler {
	return &Handler{
		service:   service,
		clicks:    clicks,
		baseURL:   baseURL,
		jwtSecret: jwtSecret,
	}
}

// Debug endpoint to check auth status
func (h *Handler) handleAuthDebug(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())

	response := map[string]interface{}{
		"authenticated": ok,
		"user_id":       userID,
		"cookies":       r.Cookies(),
		"headers": map[string]string{
			"cookie": r.Header.Get("Cookie"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.With(auth.OptionalAuth(h.jwtSecret)).Post("/api/shorten", h.handleShorten)
	r.Get("/{code}", h.handleRedirect)

	// requrie auth to get specific user urls
	r.With(auth.RequireAuth(h.jwtSecret)).Get("/api/dashboard", h.handleGetMyURLs)
	r.With(auth.RequireAuth(h.jwtSecret)).Get("/api/auth/debug", h.handleAuthDebug)
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// OptionalAuth attaches a user ID only if the request was
	// authenticated. exists == false just means "anonymous" -- not
	// an error, so no early return here.
	var userIDPtr *int64
	if userID, exists := auth.UserIDFromContext(r.Context()); exists {
		userIDPtr = &userID
	}

	// u, err := h.service.ShortenUsingBase62Encode(r.Context(), req.URL)
	u, err := h.service.ShortenByHashing(r.Context(), userIDPtr, req.URL)

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

type myURLResponse struct {
	ShortCode   string    `json:"short_code"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *Handler) handleGetMyURLs(w http.ResponseWriter, r *http.Request) {
	// RequireAuth guarantees this is present -- no need to check `ok`.
	userID, ok := auth.UserIDFromContext(r.Context())

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := h.service.GetUserURLs(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get urls", http.StatusInternalServerError)
		return
	}

	resp := make([]myURLResponse, 0, len(urls))
	for _, u := range urls {
		resp = append(resp, myURLResponse{
			ShortCode:   u.ShortCode,
			ShortURL:    h.baseURL + "/" + u.ShortCode,
			OriginalURL: u.OriginalURL,
			CreatedAt:   u.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
