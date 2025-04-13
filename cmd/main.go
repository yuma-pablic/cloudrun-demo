package main

import (
	"api/config"
	sqlc "api/infra/sqlc"
	"api/libs/metrics"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	conn := config.InitDB()
	defer config.CloseDB()

	db := sqlc.New(conn)

	metrics := metrics.NewMetrics()

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		metrics.Requests.WithLabelValues(r.URL.Path).Inc()
		// 正常ならOKレスポンスを返す
		_, err := db.Healthcheck(context.Background())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
		}
		response := map[string]string{"status": "ok"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("failed to encode response: %v", err)
		}
	})

	r.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
