package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Fry-Fr/chirpy/internal/auth"
	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/google/uuid"
)

func LoginUser(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	usr, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, usr.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !match {
		RespondWithError(w, http.StatusUnauthorized, "email password do not match")
		return
	}

	expires_in_seconds := 3600 // default 1 hour
	if params.ExpiresInSeconds > 3600 {
		expires_in_seconds = 3600
	} else if params.ExpiresInSeconds < 3600 && params.ExpiresInSeconds > 0 {
		expires_in_seconds = params.ExpiresInSeconds
	}
	token, err := auth.MakeJWT(usr.ID, auth.GetJWTSecret(), time.Duration(expires_in_seconds)*time.Second)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type resVars struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}
	payload := resVars{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
		Token:     token,
	}
	RespondWithJSON(w, http.StatusOK, payload)
}
