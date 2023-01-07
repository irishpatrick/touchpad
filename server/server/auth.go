package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"touchpad/security"
)

func NewAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHandler(w, r)
		next.ServeHTTP(w, r)
	})
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.RequestURI, "/auth/challenge") {
		return // do nothing
	} else if strings.HasSuffix(r.RequestURI, "/auth/response") {
		return // do nothing
	}

	http.Error(w, "Forbidden", http.StatusForbidden)
}

func AuthLoginChallengeHandler(w http.ResponseWriter, r *http.Request) {
	challenge := security.NewChallenge()

	data, err := json.Marshal(challenge)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}

	w.Write(data)
}

func AuthLoginResponseHandler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	bodyLen, err := r.Body.Read(body)
	if err != nil || bodyLen == 0 {
		http.Error(w, "bad request: cannot read body", http.StatusBadRequest)
		return
	}

	var challenge security.Challenge
	err = json.Unmarshal(body, &challenge)
	if err != nil {
		http.Error(w, "bad request: unable to unmarshal", http.StatusBadRequest)
	}

	if !challenge.VerifySolution() {
		http.Error(w, "forbidden", http.StatusForbidden)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "",
		Value: "",
	})
}
