package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kx0101/blog-aggregator-bootdev/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error: %s", err)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatalf("PORT environment variable not set")
	}

	mux := http.NewServeMux()
	handlers.RegisterHandlers(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: mux,
	}

	log.Printf("Server is listening on port: %v", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error: %s", err)
	}
}
