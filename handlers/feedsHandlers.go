package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kx0101/blog-aggregator-bootdev/handlers/middlewares"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

func RegisterFeedHandlers(cfg *middlewares.APIConfig, mux *http.ServeMux, dbQueries *database.Queries) {
	mux.HandleFunc("POST /v1/feeds", cfg.MiddlewareAuth(handleCreateFeed))
}

func handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User, dbQueries *database.Queries) {
	var requestBody struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Error: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id := uuid.New()
	now := time.Now()

	feed, err := dbQueries.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      requestBody.Name,
		Url:       requestBody.Url,
		UserID:    user.ID,
	})

	if err != nil {
		log.Printf("Error inserting new feed: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating feed")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, feed)
}
