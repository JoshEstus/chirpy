package main

import (
	"encoding/json"
	"log"
	"net/http"
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
		Valid bool `json:"valid"`
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

	respondWithJSON(w, http.StatusOK, validReturn{
		Valid: true,
	})
}
