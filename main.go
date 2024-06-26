package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kx0101/blog-aggregator-bootdev/handlers"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
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

	go utils.FeedWorker(dbQueries, 5*time.Second, 10)

	log.Printf("Server is listening on port: %v", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error: %s", err)
	}
}
