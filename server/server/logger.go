package server

import (
	"log"
	"net/http"
)

func NewLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n\nREQUEST: %v\n", r)
		next.ServeHTTP(w, r)
	})
}
