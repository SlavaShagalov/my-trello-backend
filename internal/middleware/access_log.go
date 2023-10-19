package middleware

import (
	"go.uber.org/zap"
	"net/http"
)

func NewAccessLog(log *zap.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("New request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("protocol", r.Proto),
				zap.String("origin", r.Header.Get("Origin")))
			handler.ServeHTTP(w, r)
		})
	}
}
