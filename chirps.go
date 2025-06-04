package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/joshestus/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) AddChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleanedChirp, err := ValidateChirpHandler(params.Body)

	if err != nil {
		log.Printf("Chirp is to long at %d", utf8.RuneCountInString(params.Body))
		respondWithError(w, http.StatusBadRequest, "Chirp is to long", err)
		return
	}

	newChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: params.UserId,
	})

	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        newChirp.ID.UUID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		},
	})
}

func ValidateChirpHandler(chirp string) (string, error) {

	const maxChirpLength = 140

	// To Long
	if utf8.RuneCountInString(chirp) > maxChirpLength {
		log.Printf("Chirp is to long at %d", utf8.RuneCountInString(chirp))
		return "", errors.New("Chirp to long")

	}

	// Filter words
	cleaned := ProfaneFilter(chirp)
	return cleaned, nil
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
