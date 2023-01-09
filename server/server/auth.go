package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"touchpad/security"
)

func NewAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldContinue := authHandler(w, r)

		if shouldContinue {
			next.ServeHTTP(w, r)
		}
	})
}

func authHandler(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasSuffix(r.RequestURI, "/api/auth/challenge") {
		return true // do nothing, continue
	} else if strings.HasSuffix(r.RequestURI, "/api/auth/response") {
		return true // do nothing, continue
	} else if !strings.HasPrefix(r.RequestURI, "/api") {
		return true
	}

	jwtCookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Printf("no token cookie")
			http.Error(w, "Forbidden", http.StatusUnauthorized)
		} else {
			log.Printf("cookie error: %v", err)
			http.Error(w, "Forbidden", http.StatusBadRequest)
		}
		return false
	}

	if !security.ValidateJwtToken(jwtCookie.Value) {
		log.Printf("invalid token")
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return false
	}

	return true
}

func AuthLoginChallengeHandler(w http.ResponseWriter, r *http.Request) {
	challenge := security.NewChallenge()
	challenge.EraseAnswer()

	data, err := json.Marshal(challenge)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	w.Write(data)
}

func AuthLoginResponseHandler(w http.ResponseWriter, r *http.Request) {
	var challenge security.Challenge
	err := json.NewDecoder(r.Body).Decode(&challenge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if !challenge.VerifySolution() {
		http.Error(w, "forbidden", http.StatusForbidden)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: security.IssueJwtToken(),
	})
}
