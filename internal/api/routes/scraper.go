package routes

import (
	"fmt"
	"net/http"

	"github.com/enzom-uy/hltb-go-scraper/internal/api/routes/helpers"
	"github.com/go-chi/chi/v5"
)

func ScraperRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {

		gameName := r.URL.Query().Get("game_name")

		if gameName == "" {
			const errorMessage string = "No valid query received. Query param is \"game_name\""
			http.Error(w, errorMessage, 400)
			fmt.Println(errorMessage)
			return
		}

		fmt.Println("Params: ", gameName)
		helpers.SendJSON(w, "hola")
	})

	return r
}
