package custom_middlewares

import (
	"net/http"

	"github.com/enzom-uy/hltb-go-scraper/internal"
)

type AuthMiddleware struct {
	config internal.Config
}

func NewAuthMiddleware(cfg *internal.Config) *AuthMiddleware {
	return &AuthMiddleware{config: *cfg}
}

func (am *AuthMiddleware) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, `{"error": "API Key required"}`, http.StatusUnauthorized)
			return
		}

		valid := false
		for _, validKey := range am.config.AllowedAPIKeys {
			if apiKey == validKey {
				valid = true
				break
			}
		}

		if !valid {
			http.Error(w, `{"error": "Invalid API Key"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
