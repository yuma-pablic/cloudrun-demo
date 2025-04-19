package main

import (
	"api/config"
	"api/handler"
	sqlc "api/infra/sqlc"
	"api/libs/logger"
	"api/libs/metrics"
	"api/libs/trace"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	appmiddleware "api/middleware"
)

func main() {
	ctx := context.Background()
	logger.InitLogger()

	tp, err := trace.InitTracer(ctx)
	if err != nil {
		slog.Error("failed to initialize tracer", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("tracer shutdown error", slog.String("error", err.Error()))
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(appmiddleware.TracingMiddleware())

	r.Use(appmiddleware.TraceIDMiddleware)

	conn := config.InitDB()
	defer config.CloseDB()

	db := sqlc.New(conn)
	metrics := metrics.NewMetrics()

	handler.RegisterPprofRoutes(r)
	handler.RegisterMetricsRoute(r)
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		metrics.Requests.WithLabelValues(r.URL.Path).Inc()

		_, err := db.Healthcheck(ctx)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.Error("healthcheck failed", slog.String("error", err.Error()))
			return
		}

		response := map[string]string{"status": "ok"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			slog.Error("failed to encode response", slog.String("error", err.Error()))
			return
		}

		slog.InfoContext(ctx, "healthcheck success")
	})

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		slog.Error("server failed to start", slog.String("error", err.Error()))
	}
}
