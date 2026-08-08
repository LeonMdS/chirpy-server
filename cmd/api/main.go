package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/LeonMdS/chirpy-server/internal/api"
	"github.com/LeonMdS/chirpy-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env variables: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("SECRET_KEY")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	const port = "8080"

	cfg := api.NewAPIConfig(database.New(dbConn), platform, secretKey)

	mux := api.NewRouter(cfg)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
