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
	"os"

	"api/custom"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:12345"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("app"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

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
	ctx := context.Background()
	logger.InitLogger()

	tp, err := initTracer(ctx)
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

	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "chi-handler")
	})

	r.Use(custom.TraceIDMiddleware)

	conn := config.InitDB()
	defer config.CloseDB()

	db := sqlc.New(conn)
	metrics := metrics.NewMetrics()

	MountPprofRoutes(r)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		metrics.Requests.WithLabelValues(r.URL.Path).Inc()

		_, err := db.Healthcheck(r.Context()) // ← 修正：ctxではなくリクエストのcontextを使う
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

		slog.InfoContext(r.Context(), "healthcheck success")
	})

	r.Handle("/metrics", promhttp.Handler())

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		slog.Error("server failed to start", slog.String("error", err.Error()))
	}
}
