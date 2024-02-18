package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wipdev-tech/blogagg/internal/database"
)

func (api *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type InUser struct {
		Name string `json:"name"`
	}

	type OutUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Name      string    `json:"name"`
	}

	inUser := InUser{}
	err := json.NewDecoder(r.Body).Decode(&inUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed request body")
	}

	dbUser, err := api.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      inUser.Name,
	})

	if err != nil {
		log.Fatal(err)
	}

	outUser := OutUser{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
	}

	respondWithJSON(w, http.StatusCreated, outUser)
}
