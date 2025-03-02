package main

import "net/http"

func (cfg *Config) healthHandler(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, struct{}{})
}
