package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User, *database.Queries)

type apiConfig struct {
	dbQueries *database.Queries
}

func RegisterHandlers(mux *http.ServeMux, dbQueries *database.Queries) {
	cfg := &apiConfig{dbQueries: dbQueries}

	mux.HandleFunc("/v1/healthz", handleHealthz)
	mux.HandleFunc("/v1/err", handleErr)

	mux.HandleFunc("POST /v1/users", func(w http.ResponseWriter, r *http.Request) {
		handleCreateUser(w, r, dbQueries)
	})
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(handleGetUser))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(handleCreateFeed))
}

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("Authorization header is missing")
			utils.RespondWithError(w, http.StatusBadRequest, "Authorization header is missing")
			return
		}

		splitToken := strings.Split(token, "ApiKey ")
		if len(splitToken) != 2 {
			log.Println("Authorization header format must be 'ApiKey <your-api-key>'")
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid Authorization header format")
			return
		}

		apiKey := splitToken[1]
		user, err := cfg.dbQueries.GetUser(r.Context(), apiKey)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		handler(w, r, user, cfg.dbQueries)
	}
}

func handleHealthz(w http.ResponseWriter, _r *http.Request) {
	responseBody := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	utils.RespondWithJSON(w, http.StatusOK, responseBody)
}

func handleErr(w http.ResponseWriter, _r *http.Request) {
	utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func handleGetUser(w http.ResponseWriter, r *http.Request, user database.User, dbQueries *database.Queries) {
	utils.RespondWithJSON(w, http.StatusOK, user)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, dbQueries *database.Queries) {
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Error: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id := uuid.New()
	now := time.Now()

	user, err := dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      requestBody.Name,
	})
	if err != nil {
		log.Printf("Error inserting new user: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
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
