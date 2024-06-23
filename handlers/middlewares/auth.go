package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User, *database.Queries)

func (cfg *APIConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
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
		user, err := cfg.DBQueries.GetUser(r.Context(), apiKey)
		if err != nil {
			log.Printf("Error: %s", err.Error())
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		handler(w, r, user, cfg.DBQueries)
	}
}
