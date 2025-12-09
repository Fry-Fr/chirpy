package handlers

import (
	"net/http"

	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/google/uuid"
)

func DeleteChirp(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId")
	chirp_id, err := uuid.Parse(chirpId)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	userId, err := AuthenticateUser(w, r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	chirp, err := cfg.DB.GetChirp(r.Context(), chirp_id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if chirp.UserID != userId {
		RespondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp")
		return
	}
	if err = cfg.DB.DeleteChirp(r.Context(), chirp_id); err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusNoContent, nil)
}
