package main

import (
	"bytes"
	"crypto/aes"
	"gitlab.com/weregoat/crypto/util"
	"testing"
)

func TestChallenge(t *testing.T) {
	blockSize := aes.BlockSize
	// Chosen plaintext 4x16 should guarantee a repeated block
	pSize := 4*blockSize
	plainText := bytes.Repeat([]byte{'A'}, pSize)
	for i := 0; i < 50; i++ {
		cipherText, mode := encrypt(plainText)
		isECB := util.HasRepeatingBlocks(cipherText, blockSize)
		for _,block := range util.Split(cipherText, blockSize) {
			t.Logf("%q", util.EncodeToBase64(block))
		}
		switch isECB {
		case true:
			t.Log("ECB mode")
		case false:
			t.Log("CBC mode")
		}
		for _,block := range util.Split(cipherText, blockSize) {
			t.Logf("%q", util.EncodeToBase64(block))
		}
		if len(cipherText)%blockSize != 0 {
			t.Errorf("cipherText is of the wrong size %d, it should be a multiple of %d", len(cipherText), blockSize)
		}
		// Since the encrypting function add random bytes, the length should be more than the original size
		if len(cipherText) <= pSize {
			t.Errorf("no random bytes were added to the plaintext, the encrypting function must be broken")
		}

		if isECB && mode != ECB {
			t.Errorf("expecting ECB mode, but failed to detect it")
		}
		if ! isECB && mode != CBC {
			t.Errorf("expecting not ECB mode, but is not CBC")
		}
	}
}
