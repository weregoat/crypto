package ch9

import (
	"gitlab.com/weregoat/crypto/pkcs7"
	"testing"
)

func TestChallenge(t *testing.T) {
	src := "YELLOW SUBMARINE"
	size := 20
	expected := "YELLOW SUBMARINE\x04\x04\x04\x04"
	solution := string(pkcs7.Pad([]byte(src), size))
	if solution != expected {
		t.Errorf("expecting %s, got %s", expected, solution)
	}

}
