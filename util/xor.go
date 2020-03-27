package util

import (
	"bytes"
	"fmt"
	"log"
	"sort"
)

type Solution struct {
	Key            []byte
	PlainText      []byte
	FrequencyScore float64
	WordScore      int
}

/*
Fixed XOR
Write a function that takes two equal-length buffers and produces their XOR combination.

If your function works properly, then when you feed it the string:

1c0111001f010100061a024b53535009181c
... after hex decoding, and when XOR'd against:

686974207468652062756c6c277320657965
... should produce:

746865206b696420646f6e277420706c6179
*/

func FixedXORBytes(a, b []byte) ([]byte, error) {
	n := len(a)
	if len(b) != n {
		err := fmt.Errorf("different buffer length %d<>%d", len(a), len(b))
		return nil, err
	}
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst, nil
}

func RepeatingXORBytes(src []byte, key []byte) []byte {
	dst := make([]byte, len(src))
	for i, j := range src {
		dst[i] = j ^ key[i%len(key)]
	}
	return dst
}

func BreakRepeatingXOR(cypherText []byte, maxKeySize, maxSolutions int) []Solution {
	var solutions []Solution
	// Guess the key sizes
	distances := GetBlockDistances(cypherText, MinBlockSize, maxKeySize)
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].Median < distances[j].Median // Using median here... because.
	})
	if len(distances) < maxSolutions {
		maxSolutions = len(distances)
	}
	// Try the most probable keys first
	for _, distance := range distances[:maxSolutions] {
		var key []byte
		// Transpose the cyphertext
		blocks := Transpose(cypherText, distance.BlockSize)
		// For each transposed block try to decrypt it with a single byte key
		for _, block := range blocks {
			topScore := 0.0
			var keyByte byte
			for i := 0; i < 256; i++ {
				b := byte(i)
				a := bytes.Repeat([]byte{b}, len(block))
				xor, err := FixedXORBytes(a, block)
				if err != nil {
					log.Fatal(err)
				}
				// Check the frequency score of the XOR'd block
				// If the key byte is right, it should amount to English letters
				lowerCase := bytes.ToLower(xor)
				score := FrequencyScore(lowerCase, FrequencyMap(Frequencies))
				// Keep the key byte resulting in the top score
				if score > topScore {
					topScore = score
					keyByte = b
				}
			}
			// Add the byte to the guessed key
			key = append(key, keyByte)
		}
		// Now we went through all the blocks and we have a possible key (*the*
		// key we think most probable)
		solution := Solution{
			Key:       key,
			PlainText: RepeatingXORBytes(cypherText, key),
		}
		solution.FrequencyScore = FrequencyScore(
			solution.PlainText, FrequencyMap(Frequencies),
		)
		solution.WordScore = WordScore(solution.PlainText, Punctuation, Dict)
		solutions = append(solutions, solution)
	}
	return solutions
}
