package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	Session string `json:"session"`
	jwt.StandardClaims
}

func IssueJwtToken() string {
	claims := JWTClaims{
		GetSessionID(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().UnixMilli() + 5*1000,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetTokenHMACKey())
	if err != nil {
		log.Panic(err)
	}
	return tokenString
}

func ValidateJwtToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return GetTokenHMACKey(), nil
	})
	if err != nil {
		log.Print(err)
		return errors.New("token parse error")
	}

	if !token.Valid {
		return errors.New("token valid check fail")
	}

	var claims JWTClaims
	jsonStr, err := json.Marshal(token.Claims)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonStr, &claims); err != nil {
		return err
	}

	if claims.Session != GetSessionID() {
		return errors.New("bad session id")
	}

	if err := claims.Valid(); err != nil {
		return err
	}

	return nil
}
