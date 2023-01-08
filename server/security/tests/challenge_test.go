package security_test

import (
	"testing"
	"touchpad/security"
)

func TestChallengeVerify(t *testing.T) {
	security.GenerateKeys()
	challenge := security.NewChallenge()
	valid := challenge.VerifySolution()
	if !valid {
		t.Error("valid challenge returns not valid")
	}
}
