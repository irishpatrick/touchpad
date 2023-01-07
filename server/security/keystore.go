package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

var jwtKey *rsa.PrivateKey
var challengeKey *rsa.PrivateKey

var keyBits int = 4096

func GenerateKeys() {
	key, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		panic(err)
	}

	jwtKey = key

	key, err = rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		panic(err)
	}

	challengeKey = key
}

func GetJWTPublicKey() crypto.PublicKey {
	return jwtKey.Public()
}

func GetJWTPrivateKey() *rsa.PrivateKey {
	return jwtKey
}

func GetChallengePublicKey() crypto.PublicKey {
	return challengeKey.Public()
}

func GetChallengePrivateKey() *rsa.PrivateKey {
	return challengeKey
}
