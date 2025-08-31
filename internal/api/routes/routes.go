package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/enzom-uy/hltb-go-scraper/internal"
	custom_middlewares "github.com/enzom-uy/hltb-go-scraper/internal/custom-middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func Setup() http.Handler {

	cfg := internal.LoadConfig()
	authMiddleware := custom_middlewares.NewAuthMiddleware(cfg)

	r := chi.NewRouter()

	// Security stuff
	r.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	r.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))

	r.Use(middleware.Logger)
	r.Use(httprate.Limit(10, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint)))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(authMiddleware.APIKeyAuth)
		r.Get("/", func(response http.ResponseWriter, r *http.Request) {
			response.Write([]byte("Index route."))
		})
		r.Mount("/scraper", ScraperRoutes())
	})

	fmt.Println("Listening to port 3333.")
	http.ListenAndServe(":3333", r)

	return r
}
