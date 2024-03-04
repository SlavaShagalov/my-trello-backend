package middleware

import (
	"net/http"
)

const MainOrigin = "https://my-trello.com"

var AllowedOrigins = map[string]struct{}{
	MainOrigin:              {},
	"http://localhost":      {},
	"http://localhost:3000": {},
	"http://127.0.0.1":      {},
	"http://127.0.0.1:3000": {},
}

func NewCors() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := MainOrigin
			gotOrigin := r.Header.Get("Origin")
			if _, allowed := AllowedOrigins[gotOrigin]; allowed {
				origin = gotOrigin
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Vary", "Origin")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
