package ch12

import (
	"encoding/base64"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/util"
)

type Oracle struct {
	key       []byte
	Error     error
	plainText []byte
}

func New(secret string) (Oracle, error) {
	o := Oracle{}
	key, err := util.RandomBytes(ecb.BlockSize)
	if err != nil {
		o.Error = err
		return o, err
	}
	plainText, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		o.Error = err
		return o, err
	}
	o.key = key
	o.plainText = plainText
	return o, nil
}

func (o *Oracle) Encrypt(chosenPlainText []byte) []byte {
	var err error
	plainText := append(chosenPlainText, o.plainText...)
	cipherText, err := ecb.Encrypt(plainText, o.key)
	if err != nil {
		o.Error = err
	}
	return cipherText
}
