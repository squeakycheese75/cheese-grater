package middleware

import (
	"log"
	"net/http"
)

// Middleware for logging incoming and outgoing URLs
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request URL: %s", r.URL)

		next.ServeHTTP(w, r)
	})
}
