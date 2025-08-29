package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GameDurations struct {
	MainStory     string `json:"main_story"`
	MainsSides    string `json:"main_sides"`
	Completionist string `json:"completionist"`
}

type SuccessPayload struct {
	Query         string        `json:"query"`
	GameTitle     string        `json:"game_title"`
	GameId        string        `json:"game_id"`
	GameDurations GameDurations `json:"game_durations"`
}

type SuccessResponse struct {
	Message string         `json:"message"`
	Data    SuccessPayload `json:"data"`
}

func ScraperRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		gameName := r.URL.Query().Get("game_name")

		if gameName == "" {
			const errorMessage string = "No valid game name received."
			http.Error(w, errorMessage, 400)
			fmt.Println(errorMessage)
			return
		}

		fmt.Println("Params: ", gameName)
		w.Write([]byte(gameName))
	})

	return r
}
