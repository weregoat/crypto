package main

import (
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/util"
	"math/rand"
	"time"
)

const ModeCBC = "CBC"
const ModeECB = "ECB"

type Oracle struct {
	BlockSize  int
	Key        []byte
	CipherText []byte
	PlainText  []byte
	IV         []byte
	Mode       string
	Error      error
}

func New(blockSize int) Oracle {
	if blockSize < 0 {
		blockSize = 0
	}
	return Oracle{BlockSize: blockSize}
}

func (o *Oracle) Encrypt(plainText []byte) error {
	key, err := util.RandomBytes(o.BlockSize)
	if err != nil {
		o.Error = err
		return err
	}
	o.Key = key
	o.PlainText = util.RandomPad(plainText, 5, 10)
	rand.Seed(time.Now().UnixNano())
	mode := rand.Intn(2)
	switch mode {
	case 0:
		iv, err := util.RandomBytes(16)
		if err != nil {
			o.Error = err
			return err
		}
		cipherText, err := cbc.Encrypt(o.PlainText, key, iv)
		if err != nil {
			o.Error = err
			return err
		}
		o.IV = iv
		o.CipherText = cipherText
		o.Mode = ModeCBC
	case 1:
		cipherText, err := ecb.Encrypt(o.PlainText, key)
		if err != nil {
			o.Error = err
			return err
		}
		o.CipherText = cipherText
		o.Mode = ModeECB
	}
	return nil
}
