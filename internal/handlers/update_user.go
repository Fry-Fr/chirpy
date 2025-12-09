package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Fry-Fr/chirpy/internal/auth"
	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/Fry-Fr/chirpy/internal/database"
	"github.com/google/uuid"
)

func UpdateUser(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	usrId, err := AuthenticateUser(w, r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := &reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !IsValidEmail(params.Email) {
		RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updateParams := database.UpdateUserParams{
		ID:             usrId,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.DB.UpdateUser(r.Context(), updateParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
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
	RespondWithJSON(w, http.StatusOK, payload)
}
