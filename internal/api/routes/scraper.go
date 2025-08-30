package routes

import (
	// "fmt"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/enzom-uy/hltb-go-scraper/internal/api/handlers"
	"github.com/enzom-uy/hltb-go-scraper/internal/api/routes/helpers"
	"github.com/go-chi/chi/v5"
)

func ScraperRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()
		gameName := r.URL.Query().Get("game_name")

		defer func() {
			if ctx.Err() == context.Canceled {
				log.Printf("Request cancelled after %v", time.Since(start))
			} else {
				log.Printf("Request completed in %v", time.Since(start))
			}
		}()

		result, err := handlers.QueryGame(ctx, gameName)

		if err != nil {
			if ctx.Err() == context.Canceled {
				log.Println("Request was cancelled by client.")
				return
			}
			http.Error(w, err.Error(), 400)
			return
		}

		helpers.SendJSONSuccess(w, result)
	})

	return r
}
