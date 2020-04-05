package ch1

import (
	"encoding/hex"
	"gitlab.com/weregoat/crypto/util"
	"testing"
)

/*

 */
func TestChallenge1(t *testing.T) {
	challenge := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	solution := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	src, err := hex.DecodeString(challenge)
	if err != nil {
		t.Error(err)
	}
	enc := util.EncodeToBase64(src)
	if string(enc) != solution {
		t.Errorf("expecting %s got %s", solution, enc)
	}
}
