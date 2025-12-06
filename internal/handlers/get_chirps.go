package handlers

import (
	"net/http"
	"time"

	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/google/uuid"
)

func GetChirps(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	type resVars struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	payload := make([]resVars, len(chirps))
	for i, chirp := range chirps {
		payload[i] = resVars{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}
	RespondWithJSON(w, http.StatusOK, payload)
}
