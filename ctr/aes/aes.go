package aes

import (
	"crypto/aes"
	"gitlab.com/weregoat/crypto/util"
)

const blockSize = aes.BlockSize

func processBlock(key, iv, src, dst []byte) error {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	out := make([]byte, blockSize)
	cipher.Encrypt(out, iv)
	xor, err := util.FixedXORBytes(out,src)
	if err != nil {
		return err
	}
	for i := 0; i < len(xor); i++ {
		dst[i] = xor[i]
	}
	return nil
}