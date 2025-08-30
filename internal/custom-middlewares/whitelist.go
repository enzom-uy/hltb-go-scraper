package custom_middlewares

import (
	"fmt"
	"net/http"
	"os"
	"slices"
)

func CORSWithWhitelist() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			fmt.Println("Origin: ", origin)

			allowedOrigins := getAllowedOrigins()

			if isOriginAllowed(origin, allowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			if !isOriginAllowed(origin, allowedOrigins) {
				http.Error(w, "Origin not allowed", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getAllowedOrigins() []string {
	allowed := []string{
		"http://localhost:3333",
	}

	if prodDomain := os.Getenv("PROD_DOMAIN"); prodDomain != "" {
		allowed = append(allowed, prodDomain)
	}

	return allowed
}

func isOriginAllowed(origin string, allowed []string) bool {
	return slices.Contains(allowed, origin)
}
