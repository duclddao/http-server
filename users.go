package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var params map[string]any
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}
	createdUser, err := cfg.dbQueries.CreateUser(r.Context(), params["email"].(string))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating a new user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, createdUser)
}
