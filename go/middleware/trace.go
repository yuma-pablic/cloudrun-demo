package appmiddleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TracingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := otelhttp.NewHandler(next, "HTTP",
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				pattern := chi.RouteContext(r.Context()).RoutePattern()
				if pattern != "" {
					return r.Method + " " + pattern
				}
				return r.Method + " " + r.URL.Path // fallback
			}),
		)
		return handler
	}
}
