// All feed-related handlers and types are here
package main

import (
	"encoding/json"
	"fmt"
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

type feedWithFollow struct {
	Feed   feed   `json:"feed"`
	Follow follow `json:"follow"`
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
		respondWithError(w, http.StatusConflict, "Duplicate!!!")
		return
	}

	dbFollow, err := api.DB.CreateFollow(r.Context(), database.CreateFollowParams{
		ID:        uuid.New(),
		FeedID:    dbFeed.ID,
		UserID:    dbFeed.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println("Couldn't create follow: ", err.Error())
	}

	payload := feedWithFollow{
		Feed:   dbFeedToFeed(dbFeed),
		Follow: dbFollowToFollow(dbFollow),
	}

	respondWithJSON(w, http.StatusCreated, payload)
}

func (api *apiConfig) handleFeedsGet(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := api.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving feeds: "+err.Error())
		return
	}

	outFeeds := []feed{}
	for _, dbf := range dbFeeds {
		outFeeds = append(outFeeds, dbFeedToFeed(dbf))
	}

	respondWithJSON(w, http.StatusOK, outFeeds)
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
