package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// ミドルウェア設定（ロギングやリカバリなど）
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// /healthcheck エンドポイント定義
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"status": "ok"}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", r)
}
