package handlers

import (
	"net/http"
	"os"

	"github.com/Fry-Fr/chirpy/internal/config"
)

func AdminReset(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.FileserverHits.Store(0)

	if err := cfg.DB.ResetUsers(r.Context()); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
