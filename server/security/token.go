package security

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
)

func IssueJwtToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session": GetSessionID(),
		"issued":  0,
		"expires": 0,
	})
	tokenString, err := token.SignedString(GetTokenHMACKey())
	if err != nil {
		log.Panic(err)
	}
	return tokenString
}

func ValidateJwtToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return GetTokenHMACKey(), nil
	})
	if err != nil {
		log.Print(err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("%v\n", claims)
	} else {
		log.Print(err)
		return false
	}

	return true
}
