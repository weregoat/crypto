package util

import (
	"math"
	"strings"
	"testing"
)

func TestWordRank(t *testing.T) {
	dict := [][]byte{[]byte("the"), []byte("be")}
	punctuation := []byte{' '}

	tests := []struct {
		Plaintext    string
		Punctuation  []byte
		Dictionary   [][]byte
		ExpectedRank int
	}{
		{"thethe", punctuation, dict, 0},
		{"the the", punctuation, dict, 1},
		{"the,the", punctuation, dict, 0},
		{"the the ", punctuation, dict, 2},
		{"thethe", []byte{}, dict, 2},
		{"be the Bee ", punctuation, dict, 2},
		{"Be The Bee", punctuation, dict, 2},
	}
	for _, test := range tests {
		rank := WordScore([]byte(test.Plaintext), test.Punctuation, test.Dictionary)
		if rank != test.ExpectedRank {
			t.Errorf("expecting rank %d , got %d", test.ExpectedRank, rank)
		}
	}
}

func TestFrequencyScore(t *testing.T) {
	tests := []string{
		"It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.",
		"MARLEY was dead: to begin with. There is no doubt whatever about that.",
		"My father’s family name being Pirrip, and my Christian name Philip, my infant tongue could make of both names nothing longer or more explicit than Pip. So, I called myself Pip, and came to be called Pip.",
		"Fifteen men on the dead man's chest\nYo-ho-ho, and a bottle of rum!",
		"To Sherlock Holmes she is always _the_ woman. I have seldom heard him mention her under any other name.",
		"Stately, plump Buck Mulligan came from the stairhead, bearing a bowl of lather on which a mirror and a razor lay crossed.",
	}
	frequencyMap := FrequencyMap(Frequencies)
	for _, test := range tests {
		plaintext := strings.ToLower(test)
		score := FrequencyScore([]byte(plaintext), frequencyMap)
		if math.Abs(score-0.002) > 0.0005 { // Because of normalisation 0.0005 is arbitrary.
			t.Errorf("expecting score for %q around 0.002, got %.5f", test, score)
		}
	}
	gibberish := strings.ToLower("IeeacdmGIyfcaokzefdnelhkied")
	score := FrequencyScore([]byte(gibberish), frequencyMap)
	if math.Abs(score-0.002) < 0.0005 {
		t.Errorf("expecting a very lower score on gibberish, got %.5f", score)
	}

}
