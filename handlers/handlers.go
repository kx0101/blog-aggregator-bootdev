package handlers

import (
	"net/http"

	"github.com/kx0101/blog-aggregator-bootdev/handlers/middlewares"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
)

func RegisterHandlers(mux *http.ServeMux, dbQueries *database.Queries) {
	cfg := &middlewares.APIConfig{DBQueries: dbQueries}

	RegisterFeedHandlers(cfg, mux, dbQueries)
	RegisterUserHandlers(cfg, mux, dbQueries)
	RegisterUtilsHandlers(cfg, mux, dbQueries)
}
