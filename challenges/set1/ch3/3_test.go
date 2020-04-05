package main

import (
	"bytes"
	"encoding/hex"
	"gitlab.com/weregoat/crypto/util"
	"log"
	"math"
	"testing"
)

/*

Single-byte XOR cipher
The hex encoded string:

1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736
... has been XOR'd against a single character. Find the key, decrypt the message.

You can do this by hand. But don't: write code to do it for you.

How? Devise some method for "scoring" a piece of English plaintext. Character frequency is a good metric. Evaluate each output and choose the one with the best score.

Achievement Unlocked
You now have our permission to make "ETAOIN SHRDLU" jokes on Twitter.
*/

func TestChallenge3(t *testing.T) {
	/* "Cooking MC's like a pound of bacon" with key 58('X') */
	cypher := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"

	decoded, err := hex.DecodeString(cypher)
	if err != nil {
		t.Error(err)
	}
	topScore := 0.0
	topScorePlaintext := ""
	var topScoreKey byte
	for i := 0; i < 256; i++ {
		b := byte(i)
		a := bytes.Repeat([]byte{b}, len(decoded))
		plainText, err := util.FixedXORBytes(a, decoded)
		if err != nil {
			log.Fatal(err)
		}
		lowercase := bytes.ToLower(plainText)
		score := util.FrequencyScore(lowercase, util.FrequencyMap(util.Frequencies))
		if score > topScore {
			topScore = score
			topScorePlaintext = string(plainText)
			topScoreKey = b
		}
	}
	if topScoreKey != byte('X') {
		t.Errorf("expecting key to be byte 'X', but got %q", topScoreKey)
	}
	if topScorePlaintext != "Cooking MC's like a pound of bacon" {
		t.Errorf("expecting highest ranken plaintex to be 'Cooking MC's like a pound of bacon', but got %+q", topScorePlaintext)
	}
	if math.Abs(topScore-0.0019369260) > 0.000001 {
		t.Errorf("expecting rank to be around 0.0019369260, but got %.10f", topScore)
	}

}
