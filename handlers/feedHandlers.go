package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kx0101/blog-aggregator-bootdev/handlers/middlewares"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_feched_at"`
}

func RegisterFeedHandlers(cfg *middlewares.APIConfig, mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/feeds", cfg.MiddlewareAuth(handleCreateFeed))
	mux.HandleFunc("GET /v1/feeds", func(w http.ResponseWriter, r *http.Request) {
		handleGetFeeds(w, r, cfg.DBQueries)
	})
}

func handleGetFeeds(w http.ResponseWriter, r *http.Request, dbQueries *database.Queries) {
	feedsFromDb, err := dbQueries.GetFeeds(r.Context())
	if err != nil {
		log.Printf("Error fetching feeds: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching feeds")
		return
	}

	var feeds []Feed
	var ids []uuid.UUID
	if feedsFromDb == nil {
		feeds = []Feed{}
	}

	for _, feed := range feedsFromDb {
		ids = append(ids, feed.ID)
		feeds = append(feeds, DatabaseFeedToFeed(feed))
	}

	err = dbQueries.UpdateLastFetchedAt(r.Context(), ids)
	if err != nil {
		log.Printf("Error updating last_fetched_at for feeds: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating last_fetched_at")
		return
	}

	log.Print("Fetched feeds")
	utils.RespondWithJSON(w, http.StatusOK, feeds)
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

	log.Printf("New feed created with id: %s", id)
	utils.RespondWithJSON(w, http.StatusCreated, feed)
}

func NullableTimeToTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}

	return nil
}

func DatabaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:            dbFeed.ID,
		CreatedAt:     dbFeed.CreatedAt,
		UpdatedAt:     dbFeed.UpdatedAt,
		Name:          dbFeed.Name,
		Url:           dbFeed.Url,
		UserID:        dbFeed.UserID,
		LastFetchedAt: NullableTimeToTime(dbFeed.LastFetchedAt),
	}
}
