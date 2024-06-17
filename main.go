package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kx0101/blog-aggregator-bootdev/handlers"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	dbURL := os.Getenv("dbURL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	dbQueries := database.New(db)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatalf("PORT environment variable not set")
	}

	mux := http.NewServeMux()
	handlers.RegisterHandlers(mux, dbQueries)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: mux,
	}

	log.Printf("Server is listening on port: %v", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error: %s", err)
	}
}
