package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wipdev-tech/blogagg/internal/database"
	"github.com/wipdev-tech/blogagg/internal/feedfetch"
)

type apiConfig struct {
	DB *database.Queries
}

func (api *apiConfig) FetchFeeds(n int32) {
	wg := &sync.WaitGroup{}
	ctx := context.Background()

	dbFeeds, err := api.DB.GetNextFeedsToFetch(ctx, n)
	if err != nil {
		log.Fatal("couldn't fetch -- ", err)
	}

	feeds := []feedfetch.RSSFeed{}
	for _, f := range dbFeeds {
		wg.Add(1)
		dbFeed := f

		go func() {
			defer wg.Done()

			feed, err := feedfetch.GetFomURL(dbFeed.Url)
			if err != nil {
				log.Fatal("error fetching ", dbFeed.Url, " -- ", err)
			}
			feeds = append(feeds, feed)

			api.DB.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
				ID:            dbFeed.ID,
				LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
			})
		}()
	}

	wg.Wait()
	fmt.Println("\nFetched:")
	for _, r := range feeds {
		fmt.Println(r.Channel.Title)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := ":" + os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DBURL")
	if dbURL == "" {
		log.Fatal("DBURL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	api := apiConfig{}
	api.DB = database.New(db)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			api.FetchFeeds(2)
		}
	}()

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	rv1 := chi.NewRouter()
	rv1.Get("/readiness", handleReady)
	rv1.Get("/error", handleError)

	rv1.Post("/users", api.handleCreateUser)
	rv1.Get("/users", api.middlewareAuth(api.handleUsersGet))

	rv1.Post("/feeds", api.middlewareAuth(api.handleFeedsCreate))
	rv1.Get("/feeds", api.handleFeedsGet)

	rv1.Post("/follows", api.middlewareAuth(api.handleFollowsCreate))
	rv1.Get("/follows", api.middlewareAuth(api.handleFollowsGet))
	rv1.Delete("/follows/{followID}", api.middlewareAuth(api.handleFollowsDelete))

	r.Mount("/v1", rv1)

	server := http.Server{Addr: port, Handler: r}
	fmt.Println("Listening at port", port)
	log.Fatal(server.ListenAndServe())
}
