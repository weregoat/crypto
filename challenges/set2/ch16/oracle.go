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
	var q string
	// to prove that the attack is not specific to the way the sanitize input
	switch util.RandomInt(0, 1) {
	case 0:
		q = strings.ReplaceAll(src, ";", "\";\"")
		q = strings.ReplaceAll(q, "=", "\"=\"")
	case 1:
		q = url.QueryEscape(src)
	}

	input := strings.Join([]string{Prefix, strings.TrimSpace(q), Postfix}, "")
	plaintext := pkcs7.Pad([]byte(input), cbc.BlockSize)
	//plainPrint(plaintext)
	ciphertext, err := cbc.Encrypt(plaintext, o.Key, o.IV)
	if err != nil {
		log.Print(err)
	}
	return ciphertext
}

func (o Oracle) IsAdmin(src []byte) bool {
	plaintext, err := cbc.Decrypt(src, o.Key, o.IV)
	//plainPrint(plaintext)
	if err != nil {
		log.Print(err)
	}
	return strings.Contains(string(plaintext), Admin)
}
