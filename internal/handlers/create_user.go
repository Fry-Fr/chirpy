package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Fry-Fr/chirpy/internal/auth"
	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/Fry-Fr/chirpy/internal/database"
	"github.com/google/uuid"
)

func CreateUser(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	email := strings.TrimSpace(params.Email)
	if email == "" {
		RespondWithError(w, http.StatusBadRequest, "email is required")
		return
	}
	if !IsValidEmail(email) {
		RespondWithError(w, http.StatusBadRequest, "invalid email")
		return
	}

	hashed_pw, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user_params := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_pw,
	}
	user, err := cfg.DB.CreateUser(r.Context(), user_params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong creating user")
		return
	}
	type resVars struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	payload := resVars{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	RespondWithJSON(w, http.StatusCreated, payload)
}
