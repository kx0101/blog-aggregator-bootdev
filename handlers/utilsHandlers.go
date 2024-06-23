package handlers

import (
	"net/http"

	"github.com/kx0101/blog-aggregator-bootdev/handlers/middlewares"
	"github.com/kx0101/blog-aggregator-bootdev/internal/database"
	"github.com/kx0101/blog-aggregator-bootdev/utils"
)

func RegisterUtilsHandlers(cfg *middlewares.APIConfig, mux *http.ServeMux, dbQueries *database.Queries) {
	mux.HandleFunc("/v1/healthz", handleHealthz)
	mux.HandleFunc("/v1/err", handleErr)
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
