package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wipdev-tech/blogagg/internal/database"
)

type follow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (api *apiConfig) handleFollowsCreate(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	inFollow := struct {
		FeedID uuid.UUID `json:"feed_id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&inFollow)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed request body")
		return
	}

	dbFollow, err := api.DB.CreateFollow(r.Context(), database.CreateFollowParams{
		ID:        uuid.New(),
		FeedID:    inFollow.FeedID,
		UserID:    dbUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusConflict, "Problem creating the follow: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dbFollowToFollow(dbFollow))
}

func (api *apiConfig) handleFollowsDelete(w http.ResponseWriter, r *http.Request, _ database.User) {
	followID := chi.URLParam(r, "followID")
	followUUID, err := uuid.Parse(followID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed follow ID: "+err.Error())
		return
	}

	err = api.DB.DeleteFollow(r.Context(), followUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func dbFollowToFollow(f database.Follow) follow {
	return follow{
		ID:        f.ID,
		FeedID:    f.FeedID,
		UserID:    f.UserID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
