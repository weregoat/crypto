package ch12

import (
	"gitlab.com/weregoat/crypto/cbc/aes"
	"testing"
)

const secret = "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"

func TestGuessBlockSize(t *testing.T) {
	oracle, err := New(secret)
	if err != nil {
		t.Error(err)
	}
	blockSize := GuessBlockSize(oracle)
	if blockSize != aes.BlockSize {
		t.Errorf("expecting blocksize of %d, got %d", aes.BlockSize, blockSize)
	}
}

func TestIsECB(t *testing.T) {
	oracle, err := New(secret)
	if err != nil {
		t.Error(err)
	}
	blockSize := GuessBlockSize(oracle)
	// We know is ECB
	if !IsECB(oracle, blockSize) {
		t.Errorf("the function should have detected that is ECB, but didn't")
	}
}
