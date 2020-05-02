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
	PlainTexts []string
	Key []byte
	IV []byte
	Plaintext string
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
	o.Key = key
	o.IV = iv
	o.PlainTexts = plainTexts
	return o, nil
}

func (o *Oracle) Encrypt() (cipherText []byte, iv []byte) {
	i := 0
	if len(o.PlainTexts) > 1 {
		i = util.RandomInt(0, len(o.PlainTexts)-1)
	}
	plaintext, err := base64.StdEncoding.DecodeString(o.PlainTexts[i])
	if err != nil {
		log.Fatal(err)
	}
	o.Plaintext = string(plaintext)
	pad := pkcs7.Pad(plaintext, cbc.BlockSize)
	cipherText, err = cbc.Encrypt(pad, o.Key, o.IV)
	if err != nil {
		log.Fatal(err)
	}
	return cipherText, o.IV
}

func (o *Oracle) CheckPadding(cipherText []byte) bool {
	decrypted, err := cbc.Decrypt(cipherText, o.Key, o.IV)

	if err != nil {
		log.Fatal(err)
	}
	return pkcs7.IsPadded(decrypted)
}