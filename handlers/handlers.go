package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

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

func RegisterHandlers(mux *http.ServeMux, dbQueries *database.Queries) {
	mux.HandleFunc("/v1/healthz", handleHealthz)
	mux.HandleFunc("/v1/err", handleErr)
	mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		handleCreateUser(w, r, dbQueries)
	})
}
