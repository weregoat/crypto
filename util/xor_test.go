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
