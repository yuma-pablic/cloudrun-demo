package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	connStr := "postgres://admin:admin@localhost:5432/postgres"
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(ctx, "SELECT 'Hello from Postgres!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	fmt.Println(greeting)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"status": "ok"}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
		}
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
