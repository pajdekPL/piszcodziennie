package main

import (
	"net/http"
)

func (cfg *Config) reset(w http.ResponseWriter, req *http.Request) {

	if cfg.env != "dev" {
		respondWithError(w, http.StatusForbidden, "method forbidden", nil)
		return
	}
	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem with deleting users", err)
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
}
