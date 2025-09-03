package main

import (
	"log"

	"github.com/enzom-uy/hltb-go-scraper/internal/api/routes"
	"github.com/enzom-uy/hltb-go-scraper/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	routes.Setup()
	_, closeDB, err := db.Init()

	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}

	defer func() {
		if err := closeDB(); err != nil {
			log.Printf("Failed to close db: %v", err)
		}
	}()

}
