package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	var reqBody map[string]any

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing request", err)
		return
	}
	val, ok := reqBody["body"].(string)
	if !ok {
		respondWithError(w, 400, "Bad request body", err)
		return
	}

	if len(val) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}
	words := strings.Split(val, "")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	val = strings.Join(words, " ")
	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": val})
}
