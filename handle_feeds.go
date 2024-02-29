// All feed-related handlers and types are here
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wipdev-tech/blogagg/internal/database"
)

type feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func (api *apiConfig) handleFeedsCreate(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	inFeed := struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}{}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&inFeed)
	if err != nil || inFeed.Name == "" || inFeed.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbFeed, err := api.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      inFeed.Name,
		Url:       inFeed.URL,
		UserID:    dbUser.ID,
	})

	if err != nil {
		log.Fatal(err)
	}

	respondWithJSON(w, http.StatusCreated, dbFeedToFeed(dbFeed))
}

func dbFeedToFeed(dbFeed database.Feed) feed {
	return feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		URL:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}
