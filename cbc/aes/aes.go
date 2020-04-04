package aes

import (
	"crypto/aes"
	"fmt"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
)

const BlockSize = aes.BlockSize

// Encrypt returns a ciphertext encrypted with AES in CBC mode.
func Encrypt(plainText, key, iv []byte) ([]byte, error) {
	IV := make([]byte, len(iv)) // Need to make a copy
	copy(IV, iv)
	if len(plainText)%BlockSize != 0 {
		plainText = pkcs7.Pad(plainText, BlockSize)
	}
	cipherText := make([]byte, len(plainText))
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return cipherText, err
	}
	i := 0
	for {
		bs := i * BlockSize
		if bs >= len(cipherText) {
			break
		}
		be := (i + 1) * BlockSize
		xor, err := util.FixedXORBytes(plainText[bs:be], IV)
		if err != nil {
			return cipherText, err
		}
		cipher.Encrypt(cipherText[bs:be], xor)
		i++
		copy(IV, cipherText[bs:be])
	}

	return cipherText, nil
}

// Decrypt returns a plaintext from a ciphertext encrypted with AES in CBC mode.
func Decrypt(cipherText, key, iv []byte) ([]byte, error) {
	IV := make([]byte, len(iv))
	copy(IV, iv)
	var plainText []byte
	if len(cipherText)%BlockSize != 0 {
		err := fmt.Errorf("ciphertext byte length(%d) must be a multiple of %d bytes", len(cipherText), BlockSize)
		return plainText, err
	}
	plainText = make([]byte, len(cipherText))
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return plainText, err
	}
	i := 0
	for {
		bs := i * BlockSize
		if bs >= len(cipherText) {
			break
		}
		be := (i + 1) * BlockSize
		decripted := make([]byte, BlockSize)
		cipher.Decrypt(decripted, cipherText[bs:be])
		xor, err := util.FixedXORBytes(decripted, IV)
		if err != nil {
			return plainText, err
		}
		i++
		copy(plainText[bs:be], xor)
		copy(IV, cipherText[bs:be])

	}
	return pkcs7.RemovePadding(plainText), nil
}
