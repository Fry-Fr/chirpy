package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

func HandleHealthzStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) HandleMetricsReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) HandleMetricsLoad(w http.ResponseWriter, r *http.Request) {
	hits := cfg.FileserverHits.Load()
	display := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(display))
}

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := reqParams{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	const max_chirp_len = 140
	if len(params.Body) > max_chirp_len {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleaned_body := profaneWordSanitizer(params.Body)

	type resVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	payload := resVals{CleanedBody: cleaned_body}
	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email string `json:"email"`
	}
	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	email := strings.TrimSpace(params.Email)
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}
	if !isValidEmail(email) {
		respondWithError(w, http.StatusBadRequest, "invalid email")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong creating user")
		return
	}
	type resVars struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	payload := resVars{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, payload)
}

// Handler Helpers ...
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func profaneWordSanitizer(s string) string {
	profane_word_list := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	split_s := strings.Split(s, " ")
	for i, word := range split_s {
		l := strings.ToLower(word)
		_, ok := profane_word_list[l]
		if ok {
			split_s[i] = "****"
		}
	}
	join_splits := strings.Join(split_s, " ")
	return join_splits
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	e := map[string]string{
		"error": msg,
	}
	dat, err := json.Marshal(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshaling data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
