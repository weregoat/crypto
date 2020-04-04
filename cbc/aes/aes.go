package aes

import (
	"fmt"
	"gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
)

const BlockSize = aes.BlockSize

// Encrypt returns a ciphertext encrypted with AES in CBC mode.
func Encrypt(plainText, key, iv []byte) ([]byte, error) {
	var cipherText []byte
	blocks := pkcs7.Split(plainText, BlockSize)
	for _, block := range blocks {
		xor, err := util.FixedXORBytes(block, iv)
		if err != nil {
			return cipherText, err
		}
		eb, err := aes.Encrypt(xor, key)
		if err != nil {
			return cipherText, err
		}
		iv = eb
		cipherText = append(cipherText, eb...)
	}
	return cipherText, nil
}

// Decrypt returns a plaintext from a ciphertext encrypted with AES in CBC mode.
func Decrypt(cipherText, key, iv []byte) ([]byte, error) {
	var plainText []byte
	if len(cipherText)%BlockSize != 0 {
		err := fmt.Errorf("ciphertext byte length(%d) must be a multiple of %d bytes", len(cipherText), BlockSize)
		return plainText, err
	}
	blocks := pkcs7.Split(cipherText, BlockSize)
	for _, b := range blocks {
		db, err := aes.Decrypt(b, key)
		if err != nil {
			return plainText, err
		}
		xor, err := util.FixedXORBytes(db, iv)
		if err != nil {
			return plainText, err
		}
		iv = b
		plainText = append(plainText, xor...)
	}
	return pkcs7.RemovePadding(plainText), nil
}
