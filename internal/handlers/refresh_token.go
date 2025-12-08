package handlers

import (
	"net/http"
	"time"

	"github.com/Fry-Fr/chirpy/internal/auth"
	"github.com/Fry-Fr/chirpy/internal/config"
)

func RefreshToken(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetRefreshToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "missing refresh token")
		return
	}
	rTkn, err := cfg.DB.GetValidRefreshToken(r.Context(), refresh_token)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	expires_in_seconds := 3600 // default 1 hour
	newAccessToken, err := auth.MakeJWT(rTkn.UserID, auth.GetJWTSecret(), time.Duration(expires_in_seconds)*time.Second)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create access token")
		return
	}

	type resVars struct {
		Token string `json:"token"`
	}

	response := resVars{
		Token: newAccessToken,
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func RevokeToken(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetRefreshToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "missing refresh token")
		return
	}
	rTkn, err := cfg.DB.GetValidRefreshToken(r.Context(), refresh_token)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	if err := cfg.DB.RevokeRefreshToken(r.Context(), rTkn.Token); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not revoke refresh token")
		return
	}
	RespondWithJSON(w, http.StatusNoContent, nil)
}
