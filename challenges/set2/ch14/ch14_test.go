package ch14

import (
	"bytes"
	"encoding/base64"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/util"
	"testing"
)

const secret = "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"


func TestGetSameBlockStart(t *testing.T) {
	blockSize := ecb.BlockSize // AES
	for i:=0; i < 20; i++ {
		key, err := util.RandomBytes(blockSize)
		if err != nil {
			t.Error(err)
		}
		prefix, err := util.RandomBytes(util.RandomInt(1,4)*blockSize)
		if err != nil {
			t.Error(err)
		}
		postfix, err := util.RandomBytes(util.RandomInt(1, 4)*blockSize)
		if err != nil {
			t.Error(err)
		}
		same, err := util.RandomBytes(blockSize)
		if err != nil {
			t.Error(err)
		}
		plaintext := append(prefix, bytes.Repeat(same, 2)...)
		plaintext = append(plaintext, postfix...)
		cipherText, err := ecb.Encrypt(plaintext, key)
		for _,j := range util.Split(cipherText, blockSize) {
			t.Logf("%x\n", j)
		}
		expected := len(prefix)
		start := GetSameBlockStart(cipherText, blockSize)
		if expected != start {
			t.Errorf("expecting identical blocks to start at byte %d, but got %d", expected, start)
		}
	}
}

func TestGetPrefixLength(t *testing.T) {
	for i:=0; i < 20; i++ {
		plainText, err := util.RandomBytes(util.RandomInt(0,4*16))
		if err != nil {
			t.Error(err)
		}
		o, err := New(base64.StdEncoding.EncodeToString(plainText))
		if err != nil {
			t.Error(err)
		}
		prefixLength := GetPrefixLength(o)
		if prefixLength != len(o.Prefix) {
			t.Errorf("expecting prefix length to be %d, got %d", len(o.Prefix), prefixLength)
		}
	}
}

func TestCPA(t *testing.T) {
	oracle, err := New(secret)
	if err != nil {
		t.Error(err)
	}
	plainText := CPA(oracle)
	solution, _ := base64.StdEncoding.DecodeString(secret)
	if ! bytes.Equal(plainText, solution) {
		t.Errorf("expecting %+q, got %+q", solution, plainText)
	}
}

func TestCh13Secret(t *testing.T) {
	// Try out this trick with the set-up in #13
	// This is the plaintext that would be encoded in the cookie
	// Where we can find the format and what we want to manipulate.
	solution := "&uid=10&role=user"
	secret := base64.StdEncoding.EncodeToString([]byte(solution))
	oracle, err := New(secret)
	if err != nil {
		t.Error(err)
	}
	plainText := CPA(oracle)
	if solution != string(plainText) {
		t.Errorf("expecting %+q, got %+q", solution, plainText)
	}
}