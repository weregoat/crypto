package main

import (
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/util"
	"log"
	"math/rand"
	"time"
)

const CBC = 0
const ECB = 1

func encrypt(plainText []byte) ([]byte, int)  { // I also return the mode for verification
	key, err := util.RandomBytes(16)
	checkError(err)
	plainText = randomPad(plainText, 5,10)
	var cipherText []byte
	rand.Seed(time.Now().UnixNano())
	mode := rand.Intn(2)
	switch mode {
	case CBC:
		iv, err := util.RandomBytes(16)
		checkError(err)
		cipherText, err = cbc.Encrypt(plainText, iv, key)
		checkError(err)
	case ECB:
		cipherText, err = ecb.Encrypt(plainText, key)
	}
	return cipherText, mode
}

// Not cryptographically secure
func randomPad(src []byte, min,max int) []byte {

	before, err := util.RandomBytes(rand.Intn(max-min)+min)
	if err != nil {
		log.Fatal(err)
	}
	after, err := util.RandomBytes(rand.Intn(max-min)+min)
	if err != nil {
		log.Fatal(err)
	}
	src = append(before, src...)
	src = append(src, after...)
	return src
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}