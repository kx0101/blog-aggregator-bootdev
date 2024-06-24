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

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func RegisterUserHandlers(cfg *middlewares.APIConfig, mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/users", func(w http.ResponseWriter, r *http.Request) {
		handleCreateUser(w, r, cfg.DBQueries)
	})
	mux.HandleFunc("GET /v1/users", cfg.MiddlewareAuth(handleGetUser))
}

func handleGetUser(w http.ResponseWriter, r *http.Request, user database.User, dbQueries *database.Queries) {
	log.Printf("Fetched user with id: %s", user.ID)
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

	log.Printf("New user created with id: %s", id)
	utils.RespondWithJSON(w, http.StatusOK, user)
}

func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}
