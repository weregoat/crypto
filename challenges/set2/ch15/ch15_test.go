package ch15

import (
	"gitlab.com/weregoat/crypto/pkcs7"
	"testing"
)

func TestChallenge(t *testing.T) {
	tests := []struct {
		Src      string
		Expected string
		Valid    bool
	}{
		{"ICE ICE BABY\x04\x04\x04\x04", "ICE ICE BABY", true},
		{"ICE ICE BABY\x05\x05\x05\x05", "ICE ICE BABY\x05\x05\x05\x05", false},
		{"ICE ICE BABY\x01\x02\x03\x04", "ICE ICE BABY\x01\x02\x03\x04", false},
	}
	for _, test := range tests {
		t.Logf("Src: %+q\n", test.Src)
		s := pkcs7.RemovePadding([]byte(test.Src))
		t.Logf("Dst: %+q\n", s)
		valid := pkcs7.IsPadded([]byte(test.Src))
		if test.Valid != valid {
			t.Errorf("expecting %t as a result to padding check, but got %t", test.Valid, valid)
		}
		if string(s) != test.Expected {
			t.Errorf("expecting %q, got %q", test.Expected, s)
		}
	}
}
