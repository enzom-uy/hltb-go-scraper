package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func buildDSNFromEnv() string {
	host := getenvDefault("DB_HOST", "localhost")
	port := getenvDefault("DB_PORT", "5432")
	user := getenvDefault("DB_USER", "enzom_uy")
	password := os.Getenv("DB_PASSWORD")
	dbname := getenvDefault("DB_NAME", "backloggd")
	sslmode := getenvDefault("DB_SSLMODE", "disable")
	timezone := getenvDefault("DB_TIMEZONE", "America/Chicago")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s timezone=%s",
		host, user, password, dbname, port, sslmode, timezone)
}

func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	// cargar .env si existe (no fatal)
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = buildDSNFromEnv()
		log.Printf("DATABASE_URL no proporcionada, usando DSN construido: %s", redactPassword(dsn))
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	outPath := "./internal/models"
	// asegurar que el directorio exista y tenga permisos de escritura
	if err := os.MkdirAll(outPath, 0o755); err != nil {
		log.Fatalf("failed to create out path %s: %v", outPath, err)
	}

	// Inicializa generador gorm.io/gen
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
	})
	g.UseDB(db)

	// Lista de tablas a generar (ajusta según tu esquema)
	tables := []string{
		"users", "collection_games", "game_genres", "games", "game_platforms",
		"game_stats", "howlongtobeat_data", "review_likes", "reviews",
		"user_activity", "user_games", "users_social_links",
	}
	for _, t := range tables {
		g.GenerateModel(t)
	}

	// LLAMADA CORRECTA: Execute no retorna valor en tu versión, solo invocarla
	g.Execute()

	log.Println("Generación completada. Revisa", outPath)
}

func redactPassword(dsn string) string {
	out := dsn
	if idx := strings.Index(out, "password="); idx != -1 {
		rest := out[idx:]
		if end := strings.Index(rest, " "); end != -1 {
			out = out[:idx] + "password=REDACTED" + rest[end:]
		} else {
			out = out[:idx] + "password=REDACTED"
		}
	}
	return out
}
