package security

import (
	"crypto"
	securerandom "crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"math/rand"
	"strconv"
)

var algorithm crypto.Hash = crypto.MD5
var challengeDigits int = 4

type Challenge struct {
	Answer             string `json:"answer"`
	AnswerDigestSigned string `json:"digest"`
}

func NewChallenge() Challenge {
	if !algorithm.Available() {
		panic("hash algorithm not available")
	}

	answer := ""
	for i := 0; i < challengeDigits; i++ {
		answer += strconv.Itoa(rand.Intn(10)) // random number [0,10)
	}

	digest, err := GetChallengePrivateKey().Sign(securerandom.Reader, []byte(answer), algorithm)
	if err != nil {
		panic(err)
	}

	return Challenge{
		Answer:             answer,
		AnswerDigestSigned: hex.EncodeToString(digest),
	}
}

func (chal *Challenge) VerifySolution() bool {
	signature, err := hex.DecodeString(chal.AnswerDigestSigned)
	if err != nil {
		panic(err)
	}

	return rsa.VerifyPKCS1v15(&GetChallengePrivateKey().PublicKey, algorithm, []byte(chal.Answer), signature) == nil
}
