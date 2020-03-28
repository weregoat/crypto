package util

import "testing"

func TestRemovePadding(t *testing.T) {

	tests := []struct {
		Src      string
		Expected string
	}{
		{"Padded\x02\x02", "Padded"},
		{"Padded\x01", "Padded"}, // Tricky, could be wrong
		{"\x01", ""},             // Empty string with one byte of padding
		{"\x02\x02", ""},         // Empty string with two bytes of padding
		{"Not really\x03\x03", "Not really\x03\x03"},
		{"This is not\x02\x02 padded either\x01\x02", "This is not\x02\x02 padded either\x01\x02"},
	}
	for _, test := range tests {
		s := RemovePadding([]byte(test.Src))
		if string(s) != test.Expected {
			t.Errorf("expecting %q, got %q", test.Expected, s)
		}
	}
}
