package ch6

import (
	"encoding/base64"
	"gitlab.com/weregoat/crypto/util"
	"io/ioutil"
	"log"
	"math"
	"testing"
)

func TestChallenge6(t *testing.T) {
	expectedKey := "Terminator X: Bring the noise"
	expectedPlaintextSample := "I'm back and I'm rin" // 20 bytes
	expectedWordScore := 94
	expectedFrequencyScore := 0.002 // +- .0005
	encoded, err := ioutil.ReadFile("6.txt")
	if err != nil {
		log.Fatal(err)
	}
	cipherText, _ := base64.StdEncoding.DecodeString(string(encoded))
	solutions := util.BreakRepeatingXOR(cipherText, 40, 1)
	if len(solutions) != 1 {
		t.Errorf("expecting 1 solution, got %d", len(solutions))
	}
	solution := solutions[0]
	if solution.WordScore != expectedWordScore {
		t.Errorf("expecting a word score of %d, got %d", expectedWordScore, solution.WordScore)
	}
	if math.Abs(solution.FrequencyScore-expectedFrequencyScore) > 0.0005 {
		t.Errorf("expecting a frequency score around %.3f, got %.10f", expectedFrequencyScore, solution.FrequencyScore)
	}
	if string(solution.Key) != expectedKey {
		t.Errorf("expecting the key to be %q, got %q", expectedKey, solution.Key)
	}
	if string(solution.PlainText[:20]) != expectedPlaintextSample {
		t.Errorf("expecting the first 20 bytes of plaintext to be %q, got %q", expectedPlaintextSample, solution.PlainText[0:20])
	}

}
