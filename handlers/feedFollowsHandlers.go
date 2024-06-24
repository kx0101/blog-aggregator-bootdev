package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kx0101/blog-aggregator-bootdev/handlers/middlewares"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

func RegisterFeedFollowsHandlers(cfg *middlewares.APIConfig, mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/feed_follows", cfg.MiddlewareAuth(handleCreateFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", cfg.MiddlewareAuth(handleGetFeedFollowsForUser))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.MiddlewareAuth(handleDeleteFeedFollow))
}

func handleGetFeedFollowsForUser(w http.ResponseWriter, r *http.Request, user database.User, dbQueries *database.Queries) {
	feeds, err := dbQueries.GetFeedFollowsForUser(r.Context(), user.ID)

	if err != nil {
		log.Printf("Error: %s getting feed follows for the user with id: %s", err.Error(), user.ID)
		utils.RespondWithError(w, http.StatusBadRequest, "Error deleting feed follow")
		return
	}

	if feeds == nil {
		feeds = []database.FeedFollow{}
	}

	log.Printf("Fetched feed follows for the user with id: %s", user.ID)
	utils.RespondWithJSON(w, http.StatusOK, feeds)
}

func handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User, dbQueries *database.Queries) {
	feedFollowID := r.PathValue("feedFollowID")
	if feedFollowID == "" {
		log.Print("No feed follow id provided")
		utils.RespondWithError(w, http.StatusBadRequest, "No feed follow id provided")
		return
	}

	id, err := uuid.Parse(feedFollowID)
	if err != nil {
		log.Println("feed follow id is not in uuid format")
		utils.RespondWithError(w, http.StatusBadRequest, "feed follow id is not in uuid format")
		return
	}

	_, err = dbQueries.GetFeedFollow(r.Context(), id)
	if err != nil {
		log.Printf("Error: %s deleting feed follow with id: %s", err.Error(), id)
		utils.RespondWithError(w, http.StatusBadRequest, "there is feed follow with that id")
		return
	}

	err = dbQueries.DeleteFeedFollow(r.Context(), id)
	if err != nil {
		log.Printf("Error: %s deleting feed follow with id: %s", err.Error(), id)
		utils.RespondWithError(w, http.StatusBadRequest, "Error deleting feed follow")
		return
	}

	responseBody := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Feed follow with id: %s has been successfully deleted", id),
	}

	log.Printf("Feed follow with id: %s has been successfully deleted", id)
	utils.RespondWithJSON(w, http.StatusOK, responseBody)
}

func handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User, dbQueries *database.Queries) {
	var requestBody struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Error: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	feed, err := dbQueries.GetFeed(r.Context(), requestBody.FeedId)
	if err != nil {
		log.Printf("Error fetching feed with id: %s", requestBody.FeedId)
		utils.RespondWithError(w, http.StatusInternalServerError, "There is no feed with this id")
		return
	}

	if feed.ID != requestBody.FeedId {
		log.Print("Feed ids do not match")
		utils.RespondWithError(w, http.StatusBadRequest, "Error fetching feed")
		return
	}

	id := uuid.New()
	now := time.Now()

	feedFollow, err := dbQueries.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        id,
		UserID:    user.ID,
		FeedID:    requestBody.FeedId,
		CreatedAt: now,
		UpdatedAt: now,
	})

	if err != nil {
		log.Printf("Error creating a new feed follow: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating a new feed follow")
		return
	}

	responseBody := struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:       feed,
		FeedFollow: feedFollow,
	}

	log.Printf("New feed follow created with id: %s", id)
	utils.RespondWithJSON(w, http.StatusCreated, responseBody)
}
