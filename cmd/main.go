package main

import (
	"api/config"
	sqlc "api/infra/sqlc"
	"api/libs/logger"
	"api/libs/metrics"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MountPprofRoutes(r chi.Router) {
	r.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
		r.Get("/block", pprof.Handler("block").ServeHTTP)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.Get("/heap", pprof.Handler("heap").ServeHTTP)
		r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
		r.Get("/profile", pprof.Profile)
		r.Post("/symbol", pprof.Symbol)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		r.Get("/trace", pprof.Trace)
	})
}
func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	conn := config.InitDB()
	defer config.CloseDB()

	db := sqlc.New(conn)

	metrics := metrics.NewMetrics()

	logger.InitLogger()

	MountPprofRoutes(r)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		metrics.Requests.WithLabelValues(r.URL.Path).Inc()
		// 正常ならOKレスポンスを返す
		_, err := db.Healthcheck(context.Background())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.Error("healthcheck failed", slog.String("error", err.Error()))
		}
		response := map[string]string{"status": "ok"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.Error("failed to encode response", slog.String("error", err.Error()))
		}
		slog.Info("healthcheck success")
	})

	r.Handle("/metrics", promhttp.Handler())
	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("server failed to start", slog.String("error", err.Error()))
	}
}
