package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/Fry-Fr/chirpy/internal/database"
	"github.com/google/uuid"
)

func PolkaWebhooks(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqVars struct {
		EventType string `json:"event"`
		Data      struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	params := &reqVars{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if params.EventType != "user.upgraded" {
		RespondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	updateUserParams := database.UpdateUserChirpyRedParams{
		ID:          params.Data.UserId,
		IsChirpyRed: true,
	}
	_, err := cfg.DB.UpdateUserChirpyRed(r.Context(), updateUserParams)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusNoContent, nil)
}
