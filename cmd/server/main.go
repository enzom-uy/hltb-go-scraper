package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	http.ListenAndServe(":3333", r)
	// if err := searchGame(gameName); err != nil {
	// 	log.Fatal(err)
	// }
}
