package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"
)

func HealthzHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func ValidateChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type validReturn struct {
		CleanedBody string `json:"cleaned_body"`
	}

	const maxChirpLength = 140

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// To Long
	if utf8.RuneCountInString(params.Body) > maxChirpLength {
		log.Printf("Chirp is to long at %d", utf8.RuneCountInString(params.Body))
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	// Filter words
	cleaned := ProfaneFilter(params.Body)

	respondWithJSON(w, http.StatusOK, validReturn{
		CleanedBody: cleaned,
	})
}

func ProfaneFilter(s string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	sSplit := strings.Split(s, " ")

	newWord := make([]string, 0)
	for _, word := range sSplit {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			newWord = append(newWord, "****")
		} else {
			newWord = append(newWord, word)
		}

	}
	return strings.Join(newWord, " ")

}
