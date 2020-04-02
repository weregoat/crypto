package aes

import (
	"gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
)

const BlockSize = aes.BlockSize

// Encrypt returns a ciphertext encrypted with AES in CBC mode.
func Encrypt(plainText, IV, key []byte) ([]byte, error) {
	var cipherText []byte
	blocks := pkcs7.Split(plainText, BlockSize)
	for _, block := range blocks {
		xor, err := util.FixedXORBytes(block, IV)
		if err != nil {
			return cipherText, err
		}
		eb, err := aes.Encrypt(xor, key)
		if err != nil {
			return cipherText, err
		}
		IV = eb
		cipherText = append(cipherText, eb...)
	}
	return cipherText, nil
}

// Decrypt returns a plaintext from a ciphertext encrypted with AES in CBC mode.
func Decrypt(cipherText, IV, key []byte) ([]byte, error) {
	var plainText []byte
	blocks := pkcs7.Split(cipherText, BlockSize)
	for _, b := range blocks {
		db, err := aes.Decrypt(b, key)
		if err != nil {
			return plainText, err
		}
		xor, err := util.FixedXORBytes(db, IV)
		if err != nil {
			return plainText, err
		}
		IV = b
		t := pkcs7.RemovePadding(xor)
		plainText = append(plainText, t...)
	}
	return plainText, nil
}
