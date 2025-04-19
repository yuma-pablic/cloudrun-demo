package handler

import (
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func RegisterPprofRoutes(r chi.Router) {
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
