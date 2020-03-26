package util

import (
	"bytes"
)

// List of frequencies generated through _FrequencyCalculator_ on a
// small corpus of English classics from Project Gutenberg (Moby Dick,
// Pride and Prejudice, etc...)
// Notice that it doesn't strictly follow the order on Wikipedia:
// https://en.wikipedia.org/wiki/Letter_frequency
// Also, probably, the space is over represented because of the formatting
// in the corpus text.
var Frequencies = []Frequency{
	{Character: ' ', Value: 0.16864},
	{Character: 'e', Value: 0.09282},
	{Character: 't', Value: 0.06780},
	{Character: 'a', Value: 0.05966},
	{Character: 'o', Value: 0.05751},
	{Character: 'n', Value: 0.05129},
	{Character: 'i', Value: 0.05083},
	{Character: 'h', Value: 0.04682},
	{Character: 's', Value: 0.04616},
	{Character: 'r', Value: 0.04320},
	{Character: 'd', Value: 0.03353},
	{Character: 'l', Value: 0.03029},
	{Character: 'u', Value: 0.02179},
	{Character: 'm', Value: 0.01943},
	{Character: 'w', Value: 0.01829},
	{Character: 'c', Value: 0.01757},
	{Character: 'f', Value: 0.01636},
	{Character: 'y', Value: 0.01594},
	{Character: 'g', Value: 0.01543},
	{Character: ',', Value: 0.01404},
	{Character: 'p', Value: 0.01248},
	{Character: 'b', Value: 0.01173},
	{Character: '.', Value: 0.00872},
	{Character: 'v', Value: 0.00699},
	{Character: 'k', Value: 0.00630},
	{Character: '\'', Value: 0.00195},
	{Character: ';', Value: 0.00183},
	{Character: 'j', Value: 0.00129},
	{Character: 'x', Value: 0.00104},
	{Character: 'q', Value: 0.00079},
	{Character: 'z', Value: 0.00040},
}

var Dict = [][]byte{
	[]byte("the"),
	[]byte("be"),
	[]byte("to"),
	[]byte("of"),
	[]byte("and"),
	[]byte("a"),
	[]byte("in"),
	[]byte("that"),
	[]byte("have"),
}
var Punctuation = []byte{' ', ',', '.', ';'}

// Match compares a plaintext with a list of words separated by punctuation and
// returns the number of matches (not normalised).
func WordScore(plainText []byte, punctuation []byte, dictionary [][]byte) int {
	if len(plainText) == 0 {
		return 0
	}
	/* always use lowercase for simplicity */
	lPlainText := bytes.ToLower(plainText)
	var words [][]byte
	if len(punctuation) == 0 {
		words = dictionary
	} else {
		for _, word := range dictionary {
			for _, p := range punctuation {
				word = append(word, p)
				words = append(words, word)
			}
		}
	}
	total := 0
	for _, word := range words {
		if len(word) > 0 {
			lWord := bytes.ToLower(word)
			total += bytes.Count(lPlainText, lWord)
		}
	}
	return total
}

// Frequency is a simple struct to store frequencies, so they can be put
// in a slice and sorted.
type Frequency struct {
	Character byte
	Value     float64
}

// FrequencyScore returns a normalised number that doesn't really mean nothing
// by itself, but it can be used to compare between plaintexts to guess which
// ones are possibly gibberish.
func FrequencyScore(plainText []byte, frequencies map[byte]float64) float64 {
	score := 0.0
	for i := 0; i < len(plainText); i++ {
		for character, frequency := range frequencies {
			if plainText[i] == character {
				score += frequency
				break
			}
		}
	}
	/*
		Because a longer text or a longer list of frequencies
		will yield a higher score than a subset of the same text (or frequencies),
		we normalise it.
	*/
	normalised := score / float64(len(plainText)*len(frequencies))
	return normalised
}

func FrequencyCalculator(source []byte, letters []byte) []Frequency {
	var frequencies []Frequency
	counts := make(map[byte]int)
	total := len(source)
	for i := 0; i < total; i++ {
		for j := 0; j < len(letters); j++ {
			if source[i] == letters[j] {
				counts[letters[j]] += 1
				break
			}
		}
	}
	for letter, count := range counts {
		f := Frequency{
			Character: letter,
			Value:     float64(count) / float64(total),
		}
		frequencies = append(frequencies, f)
	}
	return frequencies
}

func FrequencyMap(frequencies []Frequency) map[byte]float64 {
	f := make(map[byte]float64)
	for _, frequency := range frequencies {
		f[frequency.Character] = frequency.Value
	}
	return f
}
