package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kx0101/blog-aggregator-bootdev/handlers/middlewares"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

func RegisterPostsHandlers(cfg *middlewares.APIConfig, mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/posts", cfg.MiddlewareAuth(handleGetPosts))
}

func handleGetPosts(w http.ResponseWriter, r *http.Request, user database.User, dbQueries *database.Queries) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid limit value")
		return
	}

	posts, err := dbQueries.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		log.Printf("Error fetching posts for user with id: %s", user.ID)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching posts for user")
		return
	}

	log.Printf("Fetched posts for user with id: %s", user.ID)
	utils.RespondWithJSON(w, http.StatusOK, posts)
}
