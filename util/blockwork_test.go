package util

import "testing"

func TestHasRepeatingBlocks(t *testing.T) {
	tests := []struct {
		Text      string
		BlockSize int
		Expected  bool
	}{
		{"012345678910", 1, true},           // 1 and 0 are repeated
		{"012345678910", 2, false},          // No two bytes block is repeated
		{"AABBBAACBC", 2, false},            // No repetition AA != BA+AC
		{"AABBAACBCC", 2, true},             // AA repeated this time
		{"ABCAAABBBCCC", 3, false},          // is not about single byte, but blocks
		{"ABCAAABBBCCCABC", 3, true},        //
		{"ABCAAABBBCCCBCACBAACB", 3, false}, // Check order in block matters

	}
	for _, test := range tests {
		r := HasRepeatingBlocks([]byte(test.Text), test.BlockSize)
		if r != test.Expected {
			t.Errorf("expecting %t about %q to have a block of %d byte repeating, but got %t", test.Expected, test.Text, test.BlockSize, r)
		}
	}
}
