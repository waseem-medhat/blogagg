package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wipdev-tech/blogagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
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

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{}))

	rv1 := chi.NewRouter()
	rv1.Get("/readiness", handleReady)
	rv1.Get("/error", handleError)
	rv1.Post("/users", api.handleCreateUser)
	rv1.Get("/users", api.middlewareAuth(api.handleUsersGet))
	rv1.Post("/feeds", api.middlewareAuth(api.handleFeedsCreate))

	r.Mount("/v1", rv1)

	server := http.Server{Addr: port, Handler: r}
	fmt.Println("Listening at port", port)
	server.ListenAndServe()
}
