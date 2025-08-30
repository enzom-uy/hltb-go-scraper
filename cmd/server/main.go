package main

import (
	"log"

	"github.com/enzom-uy/hltb-go-scraper/internal/api/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	routes.Setup()
}
