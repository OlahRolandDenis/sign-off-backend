package api

import (
	"Desktop/signoff/internal/auth"
	"Desktop/signoff/internal/db"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const AgencyIDKey contextKey = "agencyID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "" {
			id, err := db.ValidateAPIKey(r.Context(), apiKey)
			if err == nil {
				ctx := context.WithValue(r.Context(), AgencyIDKey, id)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			id, ok := auth.ValidateJWT(token)
			if ok == nil {
				ctx := context.WithValue(r.Context(), AgencyIDKey, id)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "Authentication required", http.StatusUnauthorized)
	})
}
