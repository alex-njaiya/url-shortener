package auth

import (
	"context"
	"net/http"
)

type contextKey string

const userIDContextKey contextKey = "userID"


func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			userID, err := ParseToken(cookie.Value, secret)
			if err != nil {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	return userID, ok
}