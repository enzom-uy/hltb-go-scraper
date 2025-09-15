package main

import (
	"log"
	"os"

	"github.com/enzom-uy/hltb-go-scraper/internal/api/routes"
	"github.com/enzom-uy/hltb-go-scraper/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		} else {
			log.Println("Loaded .env file for local development")
		}
	} else {
		log.Println("Running in Railway - using environment variables")
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
