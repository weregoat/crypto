package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"gitlab.com/weregoat/crypto/util"
	"log"
	"os"
	"testing"
)

/*

Detect single-character XOR
One of the 60-character strings in this file has been encrypted by single-character XOR.

Find it.

(Your code from #3 should help.)
*/

func TestChallenge4(t *testing.T) {

	solutionPlaintext := "Now that the party is jumping\n"
	solutionLine := "7b5a4215415d544115415d5015455447414c155c46155f4058455c5b523f"
	file, err := os.Open("./4.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var topRank float64
	plainText := ""
	line := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := scanner.Text()
		decoded, err := hex.DecodeString(content)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 256; i++ {
			b := byte(i)
			a := bytes.Repeat([]byte{b}, len(decoded))
			xor, err := util.FixedXORBytes(a, decoded)
			if err != nil {
				log.Fatal(err)
			}
			lowerCase := bytes.ToLower(xor)
			rank := util.FrequencyScore(lowerCase, util.FrequencyMap(util.Frequencies))
			if rank > topRank {
				topRank = rank
				plainText = string(xor)
				line = content
			}
		}
	}
	if plainText != solutionPlaintext {
		t.Errorf("expecting %+q as plaintext, got %+q", solutionPlaintext, plainText)
	}
	if solutionLine != line {
		t.Errorf("expecting %+q as line, got +%q", solutionLine, line)
	}

}
