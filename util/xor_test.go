package util

import (
	"encoding/hex"
	"testing"
)

func TestFixedXORBytes(t *testing.T) {
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
	s, err := FixedXORBytes(a, b)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < len(s); i++ {
		if s[i] != e[i] {
			t.Errorf("expecting %q at byte %d but got %q", e[i], i, s[i])
		}
	}
}

func TestFixedXORBytesError(t *testing.T) {
	a := []byte{'a', 'b'}
	b := []byte{'c', 'd', 'e'}
	_, err := FixedXORBytes(a, b)
	if err == nil {
		t.Errorf("expecting error because of different buffer lenght got nothing")
	}
}

func TestRepeatingXORBytes(t *testing.T) {
	plainText := []byte("Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal")
	expected, err := hex.DecodeString(  "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f")
	if err != nil {
		t.Error(err)
	}
	key := []byte("ICE")
	cypherText := RepeatingXORBytes(plainText, key)
	for i,j := range cypherText {
		if j != expected[i] {
			t.Errorf("expecting byte %d to be %q, but got %q", i, expected[i], j)
		}
	}

}
