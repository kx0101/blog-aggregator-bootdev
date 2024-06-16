package handlers

import (
	"github.com/kx0101/blog-aggregator-bootdev/utils"
	"net/http"
)

func handleHealthz(w http.ResponseWriter, _r *http.Request) {
	responseBody := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	utils.RespondWithJSON(w, 200, responseBody)
}

func handleErr(w http.ResponseWriter, _r *http.Request) {
	utils.RespondWithError(w, 500, "Internal Server Error")
}

func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/v1/healthz", handleHealthz)
	mux.HandleFunc("/v1/err", handleErr)
}
