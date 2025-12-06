package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strings"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ProfaneWordSanitizer(s string) string {
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

func RespondWithError(w http.ResponseWriter, code int, msg string) {
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

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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
