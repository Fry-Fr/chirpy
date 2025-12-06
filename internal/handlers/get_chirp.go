package handlers

import (
	"net/http"
	"time"

	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/google/uuid"
)

func GetChirp(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId")
	chirp_id, err := uuid.Parse(chirpId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := cfg.DB.GetChirp(r.Context(), chirp_id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusNotFound, "Not found")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type resVars struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	payload := resVars{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	RespondWithJSON(w, http.StatusOK, payload)
}
