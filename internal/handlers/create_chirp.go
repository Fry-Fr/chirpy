package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/Fry-Fr/chirpy/internal/database"
	"github.com/google/uuid"
)

func CreateChirp(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	if err := AuthenticateUser(w, r); err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	type reqParams struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	params := &reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	const max_chirp_len = 140
	if len(params.Body) > max_chirp_len {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleaned_body := ProfaneWordSanitizer(params.Body)

	chirp := database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: params.UserId,
	}
	c, err := cfg.DB.CreateChirp(r.Context(), chirp)
	if err != nil {
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
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
	RespondWithJSON(w, http.StatusCreated, payload)
}
