package server

import (
	"log"
	"net/http"
)

func NewLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n\nREQUEST: %v\nCookies: %v\n", r, r.Cookies())
		next.ServeHTTP(w, r)
	})
}
