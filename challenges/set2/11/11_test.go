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
	pSize := 4 * blockSize
	plainText := bytes.Repeat([]byte{'A'}, pSize)
	for i := 0; i < 50; i++ {
		o := New(blockSize)
		err := o.Encrypt(plainText)
		if err != nil {
			t.Error(err)
		}
		blocks := util.Split(o.CipherText, o.BlockSize)
		for _, block := range blocks {
			t.Logf("% x", block)
		}
		isECB := util.HasRepeatingBlocks(o.CipherText, o.BlockSize)
		switch isECB {
		case true:
			t.Log("ECB mode")
		case false:
			t.Log("CBC mode")
		}
		if len(o.CipherText)%blockSize != 0 {
			t.Errorf("cipherText is of the wrong size %d, it should be a multiple of %d", len(o.CipherText), blockSize)
		}
		// Since the encrypting function add random bytes, the length should be more than the original size
		if len(o.CipherText) <= pSize {
			t.Errorf("no random bytes were added to the plaintext, the encrypting function must be broken")
		}

		if isECB && o.Mode != ModeECB {
			t.Log(o.Mode)
			t.Errorf("expecting ECB mode, but failed to detect it")
		}
		if !isECB && o.Mode != ModeCBC {
			t.Log(o.Mode)
			t.Errorf("expecting not ECB mode, but is not CBC")
		}
	}
}
