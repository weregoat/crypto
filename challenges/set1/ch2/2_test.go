package ch2

import (
	"encoding/hex"
	"gitlab.com/weregoat/crypto/util"
	"testing"
)

/*
Fixed XOR
Write a function that takes two equal-length buffers and produces their XOR combination.

If your function works properly, then when you feed it the string:

1c0111001f010100061a024b53535009181c
... after hex decoding, and when XOR'd against:

686974207468652062756c6c277320657965
... should produce:

746865206b696420646f6e277420706c6179
*/
func TestChallenge2(t *testing.T) {
	a, err := hex.DecodeString("1c0111001f010100061a024b53535009181c")
	if err != nil {
		t.Error(err)
	}
	b, err := hex.DecodeString("686974207468652062756c6c277320657965")
	if err != nil {
		t.Error(err)
	}
	e, err := hex.DecodeString("746865206b696420646f6e277420706c6179")
	if err != nil {
		t.Error(err)
	}
	s, err := util.FixedXORBytes(a, b)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(s); i++ {
		if s[i] != e[i] {
			t.Errorf("expecting %q at byte %d but got %q", e[i], i, s[i])
		}
	}
}
