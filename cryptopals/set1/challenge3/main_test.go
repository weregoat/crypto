package main

import (
	"encoding/hex"
	"testing"
)

func TestChallege(t *testing.T) {
	var dict = []string{"the", "be", "to", "of", "and", "a", "in", "that", "have"}
	var punctuation = []byte{' ', ',', '.', ';'}
	/* "Cooking MC's like a pound of bacon" with key 58('X') */
	cypher := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"

	decoded, err := hex.DecodeString(cypher)
	if err != nil {
		t.Error(err)
	}
	words := make([][]byte, len(dict))
	for k, v := range dict {
		words[k] = []byte(v)
	}
	key, plainText, rank := rank(decoded, punctuation, words)
	if key != byte('X') {
		t.Errorf("expecting key to be byte 'X', but got %q", key)
	}
	if plainText != "Cooking MC's like a pound of bacon" {
		t.Errorf("expecting highest ranken plaintex to be 'Cooking MC's like a pound of bacon', but got '%s", plainText)
	}
	if rank != 2 {
		t.Errorf("expecting rank to be 2, but got %d", rank)
	}


}
