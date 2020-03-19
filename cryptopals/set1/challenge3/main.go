package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"gitlab.com/weregoat/crypto/util"
	"log"
	"os"
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

/* https://en.wikipedia.org/wiki/Most_common_words_in_English */
var dict = []string{"the", "be", "to", "of", "and", "a", "in", "that", "have"}
var puntuation = []byte{' ', ',', '.', ';'}

func main() {
	/*
			https://en.wikipedia.org/wiki/Letter_frequency
		The suggestion to use letter frequency is a good one, in a way. But...
		when I tried to apply it it showed some issues.
		For example given the above Cyphertext I could not have the plain text,
		which is 'Cooking MC's like a pound of bacon', score the higher, no matter
		various approaches and tweaks.

		So, I got thinking on and decided to switch to words matching. Given a list
		of words (possibly the most frequent ones), score based on the numbers of them
		presents.
		https://en.wikipedia.org/wiki/Most_common_words_in_English

	*/

	flag.Usage = usage
	encoded := flag.String("cypher", "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736", "hex encoded string to decypher")
	flag.Parse()
	if len(*encoded) == 0 {
		flag.Usage()
		log.Fatal("A cypher text is required")
	}

	cypherText, err := hex.DecodeString(*encoded)
	if err != nil {
		log.Fatal(err)
	}

	/* We convert the dict to an array of bytes */
	words := make([][]byte, len(dict))
	for k, v := range dict {
		words[k] = []byte(v)
	}
	var key byte
	var solution string
	max := 0
	// var solutions []Try
	for i := 0; i < 256; i++ {
		b := byte(i)
		a := bytes.Repeat([]byte{b}, len(cypherText))
		plainText, err := util.FixedXORBytes(a, cypherText)
		if err != nil {
			log.Fatal(err)
		}
		score := util.Match(plainText, puntuation, words)
		if score > max {
			max = score
			key = b
			solution = string(plainText)
		}
	}
	if max > 0 {
		fmt.Printf("Encoded cyphertext: %s\n", *encoded)
		fmt.Printf("With key %s(%q) becomes '%s' (rank %d)\n", hex.EncodeToString([]byte{key}), key, solution, max)
	}

}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s -cypher string\n", os.Args[0])
	flag.PrintDefaults()
}
