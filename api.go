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
	type genericError struct {
		Error string `json:"error"`
	}
	const maxChirpLength = 140

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respBody := genericError{
			Error: "Something went wrong",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	// To Long
	if utf8.RuneCountInString(params.Body) > maxChirpLength {
		respBody := genericError{
			Error: "Chirp is too long",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		log.Printf("Chirp is to long at %d", utf8.RuneCountInString(params.Body))
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Write(dat)
		return
	}

	type validReturn struct {
		Valid bool `json:"valid"`
	}
	respBody := validReturn{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	// Valid Chirp
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}
