package ch17

import (
	"encoding/base64"
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"log"
)

var defaultPlaintexts = []string {
"MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=",
"MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=",
"MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==",
"MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==",
"MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl",
"MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==",
"MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==",
"MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=",
"MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=",
"MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93",
}

type Oracle struct {
	plainTexts []string
	key []byte
	iv []byte
}

func NewOracle(plainTexts ...string) (Oracle, error) {
	if len(plainTexts) == 0 {
		plainTexts = defaultPlaintexts
	}
	o := Oracle{}
	key, err := util.RandomBytes(cbc.BlockSize)
	if err != nil {
		return o, err
	}
	iv, err := util.RandomBytes(cbc.BlockSize)
	if err != nil {
		return o, err
	}
	o.key = key
	o.iv = iv
	o.plainTexts = plainTexts
	return o, nil
}

func (o Oracle) Encrypt() (cipherText []byte, iv []byte) {
	i := util.RandomInt(0, len(o.plainTexts)-1)
	plaintext, err := base64.StdEncoding.DecodeString(o.plainTexts[i])
	if err != nil {
		log.Fatal(err)
	}
	pad := pkcs7.Pad(plaintext, cbc.BlockSize)
	cipherText, err = cbc.Encrypt(pad, o.key, o.iv)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (o Oracle) CheckPadding(cipherText []byte) bool {
	decrypted, err := cbc.Decrypt(cipherText, o.key, o.iv)
	if err != nil {
		log.Fatal(err)
	}
	return pkcs7.IsPadded(decrypted)
}