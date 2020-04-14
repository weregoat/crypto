package ch15

import (
	"gitlab.com/weregoat/crypto/pkcs7"
	"testing"
)

func TestChallenge(t *testing.T) {
	tests := []struct {
		Src      string
		Expected string
		Error    bool
	}{
		{"ICE ICE BABY\x04\x04\x04\x04", "ICE ICE BABY", false},
		{"ICE ICE BABY\x05\x05\x05\x05", "ICE ICE BABY\x05\x05\x05\x05", true},
		{"ICE ICE BABY\x01\x02\x03\x04", "ICE ICE BABY\x01\x02\x03\x04", true},
	}
	for _, test := range tests {
		t.Logf("Src: %+q\n", test.Src)
		s, err := pkcs7.RemovePadding([]byte(test.Src))
		t.Logf("Dst: %+q\n", s)
		if test.Error && err == nil {
			t.Error("expecting error, but got nothing")
		}
		if string(s) != test.Expected {
			t.Errorf("expecting %q, got %q", test.Expected, s)
		}
	}
}
