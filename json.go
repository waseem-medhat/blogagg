package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	errPayload := struct {
		Error string `json:"error"`
	}{Error: msg}
	respondWithJSON(w, code, errPayload)
}
