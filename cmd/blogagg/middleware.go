package main

import (
	"net/http"
	"strings"

	"github.com/wipdev-tech/blogagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (api *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, hasPrefix := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
		if !hasPrefix {
			respondWithError(w, http.StatusBadRequest, "Malformed authorization header")
			return
		}

		dbUser, err := api.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r, dbUser)
	}
}
