package ch16

import (
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"log"
	"net/url"
	"strings"
)

const Prefix = "comment1=cooking%20MCs;userdata="
const Postfix = ";comment2=%20like%20a%20pound%20of%20bacon"
const Admin = ";admin=true;"

type Oracle struct {
	Key []byte
	IV  []byte
}

func New() (Oracle, error) {
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
	return o, nil
}

func (o Oracle) Encrypt(src string) []byte {
	input := strings.Join([]string{Prefix, strings.TrimSpace(strings.ToLower(src)), Postfix}, "")

	plaintext := pkcs7.Pad([]byte(url.QueryEscape(input)), cbc.BlockSize)
	ciphertext, err := cbc.Encrypt(plaintext, o.Key, o.IV)
	if err != nil {
		log.Print(err)
	}
	return ciphertext
}

func (o Oracle) IsAdmin(src []byte) bool {
	plaintext, err := cbc.Decrypt(src, o.Key, o.IV)
	if err != nil {
		log.Print(err)
	}
	return strings.Contains(string(plaintext), Admin)
}
