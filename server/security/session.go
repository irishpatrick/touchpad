package security

import (
	"crypto/rand"
	"encoding/hex"
)

var sessionID string

func GenerateSessionID() {
	buffer := make([]byte, 256/8)
	n, err := rand.Reader.Read(buffer)
	if err != nil || n != len(buffer) {
		panic(err)
	}

	sessionID = hex.EncodeToString(buffer)
}

func GetSessionID() string {
	return sessionID
}
