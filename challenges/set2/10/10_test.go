package main

import (
	"encoding/base64"
	"gitlab.com/weregoat/crypto/cbc/aes"
	"io/ioutil"
	"log"
	"testing"
)

func TestChallenge8(t *testing.T) {
	expected := "I'm back and I'm ringin' the bell"
	encoded, err := ioutil.ReadFile("10.txt")
	if err != nil {
		log.Fatal(err)
	}
	cipherText, _ := base64.StdEncoding.DecodeString(string(encoded))
	key := []byte("YELLOW SUBMARINE")
	iv := make([]byte, 16) // Initialised a \x00
	plainText, err := aes.Decrypt(cipherText, key, iv)
	if err != nil {
		t.Error(err)
	}
	if string(plainText[0:20]) != expected[0:20] {
		t.Errorf("expecting %q, got %q", expected[0:20], plainText[0:50])

	}
}
