package main

import (
	"encoding/base64"
	"gitlab.com/weregoat/crypto/ecb/aes"
	"io/ioutil"
	"log"
	"testing"
)

func TestChalleng7(t *testing.T) {
	key := []byte("YELLOW SUBMARINE")
	expected := "I'm back and I'm rin"
	encoded, err := ioutil.ReadFile("7.txt")
	if err != nil {
		log.Fatal(err)
	}
	cipherText, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		t.Error(err)
	}

	plainText, err := aes.Decrypt(cipherText, key)
	if err != nil {
		t.Error(err)
	}

	if string(plainText[0:20]) != expected {
		t.Errorf("expecting first 20 bytes of plaintext to be %q, got %q", plainText[:20], expected)
	}

}
