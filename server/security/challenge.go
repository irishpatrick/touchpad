package security

import (
	"crypto"
	securerandom "crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var algorithm crypto.Hash = crypto.SHA256
var challengeDigits int = 4

type Challenge struct {
	Answer             string `json:"answer"`
	AnswerDigestSigned string `json:"digest"`
}

func NewChallenge() Challenge {
	rand.Seed(time.Now().UnixNano())

	if !algorithm.Available() {
		log.Panic("hash algorithm not available")
	}

	answer := ""
	for i := 0; i < challengeDigits; i++ {
		answer += strconv.Itoa(rand.Intn(10)) // random number [0,10)
	}

	hash := algorithm.New()
	hash.Write([]byte(answer))
	digest := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(securerandom.Reader, GetChallengePrivateKey(), algorithm, digest)
	if err != nil {
		log.Panic(err)
	}

	return Challenge{
		Answer:             answer,
		AnswerDigestSigned: hex.EncodeToString(signature),
	}
}

func (chal *Challenge) VerifySolution() bool {
	signature, err := hex.DecodeString(chal.AnswerDigestSigned)
	if err != nil {
		panic(err)
	}

	hash := algorithm.New()
	hash.Write([]byte(chal.Answer))
	digest := hash.Sum(nil)

	err = rsa.VerifyPKCS1v15(&GetChallengePrivateKey().PublicKey, algorithm, digest, signature)
	return err == nil
}
