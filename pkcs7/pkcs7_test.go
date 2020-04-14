package pkcs7

import (
	"testing"
)

func TestPad(t *testing.T) {
	tests := []struct {
		Src      string
		Size     int
		Expected string
	}{
		{"123", 4, "123\x01"},
		{"123", 3, "123\x03\x03\x03"},
		{"", 0, ""},
		{"1234", 0, "1234"},
		{"1234", -3, "1234"},
		{"12345678901", 5, "12345678901\x04\x04\x04\x04"},
	}
	for _, test := range tests {
		src := []byte(test.Src)
		padded := string(Pad(src, test.Size))
		if padded != test.Expected {
			t.Errorf("expecting %q for %d padding of %q, got %q", test.Expected, test.Size, test.Src, padded)
		}
	}

}

func TestRemovePadding(t *testing.T) {

	tests := []struct {
		Src      string
		Expected string
		Error    bool
	}{
		{"Padded\x02\x02", "Padded", false},
		{"Padded\x01", "Padded", false }, // Tricky, could be wrong
		{"\x01", "",false},             // Empty string with one byte of padding
		{"\x02\x02", "", false},         // Empty string with two bytes of padding
		{"Not really\x03\x03", "Not really\x03\x03",true},
		{"This is not\x02\x02 padded either\x01\x02", "This is not\x02\x02 padded either\x01\x02", false},
		{"Padded\x06\x06\x06\x06\x06\x06", "Padded", false}, // Same as first, but bigger block
		{"ICE ICE BABY\x04\x04\x04\x04", "ICE ICE BABY", false},
		{"ICE ICE BABY\x05\x05\x05\x05", "ICE ICE BABY\x05\x05\x05\x05", true},
		{"ICE ICE BABY\x01\x02\x03\x04", "ICE ICE BABY\x01\x02\x03\x04", true},
		{"ICE ICE BABY\x18", "ICE ICE BABY\x18", true},
	}
	for _, test := range tests {
		t.Logf("Src: %+q\n", test.Src)
		s, err := RemovePadding([]byte(test.Src))
		t.Logf("Dst: %+q\n", s)
		if test.Error && err == nil {
			t.Error("expecting error, but got nothing")
		}
		if string(s) != test.Expected {
			t.Errorf("expecting %q, got %q", test.Expected, s)
		}
	}
}
