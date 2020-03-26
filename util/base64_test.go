package util

import (
	"encoding/hex"
	"testing"
)

func TestEncodeToBase64(t *testing.T) {
	challenge := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	solution := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	src, err := hex.DecodeString(challenge)
	if err != nil {
		t.Error(err)
	}
	/* Check hex decoding */
	back := hex.EncodeToString(src)
	if back != challenge {
		t.Errorf("the converted byte stream %s is not the expected one %s", back, challenge)
	}
	expected := []byte(solution)
	enc := EncodeToBase64(src)
	if string(enc) != string(expected) {
		t.Errorf("expecting %s got %s", expected, enc)
	}
	for i, b := range enc {
		s := expected[i]
		if b != s {
			t.Errorf("expecting %q at byte %d got %q", s, i, b)
		}
	}
}
