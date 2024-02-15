package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	rv1 := chi.NewRouter()
	rv1.Get("/readiness", handleReady)
	rv1.Get("/error", handleError)

	r.Mount("/v1", rv1)

	port := ":" + os.Getenv("PORT")
	server := http.Server{Addr: port, Handler: r}
	fmt.Println("Listening at port", port)
	server.ListenAndServe()
}

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

func handleReady(w http.ResponseWriter, r *http.Request) {
	okPayload := struct {
		Status string `json:"status"`
	}{Status: "ok"}

	respondWithJSON(w, http.StatusOK, okPayload)
}

func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
