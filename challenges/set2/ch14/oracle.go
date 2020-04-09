package ch14

import (
	"encoding/base64"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/util"
)

type Oracle struct {
	Key       []byte
	Prefix    []byte
	Error     error
	Secret []byte
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
	length := util.RandomInt(0,48)
	prefix, err := util.RandomBytes(length)
	if err != nil {
		o.Error = err
		return o, err
	}
	o.Prefix = prefix
	o.Key = key
	o.Secret = plainText
	return o, nil
}

func (o *Oracle) Encrypt(chosenPlainText []byte) []byte {
	var err error
	prefixed := append(o.Prefix, chosenPlainText...)
	plainText := append(prefixed, o.Secret...)
	cipherText, err := ecb.Encrypt(plainText, o.Key)
	if err != nil {
		o.Error = err
	}
	return cipherText
}
