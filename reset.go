package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	// reset page counter
	cfg.fileserverHits.Swap(0)

	// empty users table
	err := cfg.dbQueries.EmptyUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error in emptying `users` table", err)
		return
	}
	//empty chirps table
	err = cfg.dbQueries.EmptyChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error in emptying `chirps` table", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
