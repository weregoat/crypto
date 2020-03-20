package util

import "testing"

func TestMatch(t *testing.T) {
	dict := [][]byte{[]byte("the"), []byte("be")}
	punctuation := []byte{' '}

	tests := []struct {
		Plaintext string
		Punctuation []byte
		Dictionary [][]byte
		ExpectedRank int
	} {
		{"thethe", punctuation, dict, 0},
		{"the the", punctuation, dict, 1},
		{"the,the", punctuation, dict, 0},
		{"the the ", punctuation, dict, 2},
		{"thethe", []byte{}, dict, 2},
		{"be the Bee ", punctuation, dict, 2},
		{"Be The Bee", punctuation, dict, 2},
	}
	for _,test := range tests {
		rank := Match([]byte(test.Plaintext), test.Punctuation, test.Dictionary)
		if rank != test.ExpectedRank {
			t.Errorf("expecting rank %d , got %d", test.ExpectedRank, rank)
		}
	}
}