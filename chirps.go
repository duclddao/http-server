package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"http-server/internal/database"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const maxChirpLength = 140

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	var reqBody map[string]any

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing request", err)
		return
	}
	body, ok := reqBody["body"].(string)
	if !ok {
		respondWithError(w, 400, "Bad request body", err)
		return
	}
	body, err = validateChirp(body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	userId, ok := reqBody["user_id"].(string)
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Bad request body", fmt.Errorf("Bad request body"))
		return
	}
	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad user_id", fmt.Errorf("Bad user_id: %w", err))
		return
	}
	createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   body,
		UserID: userIdParsed,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error in creating a new chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, createdChirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID != "" {
		parsedID, err := uuid.Parse(chirpID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error in parsing chirpID", err)
			return
		}
		chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), parsedID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "Chirp with ID "+chirpID, err)
			} else {
				respondWithError(w, http.StatusInternalServerError, "Error in getting the chirp with ID "+chirpID, err)
			}
			return
		}
		respondWithJSON(w, http.StatusOK, chirp)
	} else {
		chirps, err := cfg.dbQueries.GetChirpsOrderByCreatedAt(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting chirps", err)
			return
		}
		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func validateChirp(body string) (string, error) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	body = strings.Join(words, " ")
	return body, nil
}
