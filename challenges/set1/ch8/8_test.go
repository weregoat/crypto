package ch8

import (
	"bufio"
	"encoding/hex"
	"gitlab.com/weregoat/crypto/util"
	"os"
	"testing"
)

func TestChallenge8(t *testing.T) {
	expectedText := "d880619740a8a19b7840a8a31c810a3d08649af70dc06f4fd5d2d69c744cd283e2dd052f6b641dbf9d11b0348542bb5708649af70dc06f4fd5d2d69c744cd2839475c9dfdbc1d46597949d9c7e82bf5a08649af70dc06f4fd5d2d69c744cd28397a93eab8d6aecd566489154789a6b0308649af70dc06f4fd5d2d69c744cd283d403180c98c8f6db1f2a3f9c4040deb0ab51b29933f2c123c58386b06fba186a"
	expectedLineNumber := 133
	file, err := os.Open("8.txt")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	line := 0
	encoded := ""
	for scanner.Scan() {
		encoded = scanner.Text()
		line++ // First is 1!
		cipherText, err := hex.DecodeString(encoded)
		if err != nil {
			t.Error(err)
		}
		if util.HasRepeatingBlocks(cipherText, 16) {
			split := util.LazySplit(cipherText, 16)
			t.Logf("line %d has repeating blocks", line)
			for i, b := range split {
				t.Logf("block %d: %x\n", i, b)
			}
			break
		}
	}
	if encoded != expectedText {
		t.Errorf("expected line to be %q, but got %q", expectedText, encoded)
	}
	if line != expectedLineNumber {
		t.Errorf("expected ciphertext to be at linet %d, but I got %d", expectedLineNumber, line)
	}
}
