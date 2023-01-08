package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

var challengeKey *rsa.PrivateKey
var tokenHMACKey []byte

var keyBits int = 4096
var symmetricKeyBits int = 256

func GenerateKeys() {
	key, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		panic(err)
	}
	challengeKey = key

	tokenHMACKey := make([]byte, symmetricKeyBits/8)
	rand.Reader.Read(tokenHMACKey)
}

func GetTokenHMACKey() []byte {
	return tokenHMACKey
}

func GetChallengePublicKey() crypto.PublicKey {
	return challengeKey.Public()
}

func GetChallengePrivateKey() *rsa.PrivateKey {
	return challengeKey
}
