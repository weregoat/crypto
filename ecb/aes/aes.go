package aes

import (
	"crypto/aes"
	"fmt"
	"gitlab.com/weregoat/crypto/pkcs7"
)

const BlockSize = aes.BlockSize

/* https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#ECB */

func Decrypt(cipherText, key []byte) (plainText []byte, err error) {
	if len(cipherText)%BlockSize != 0 {
		err = fmt.Errorf("cipherText length must be a multiple of %d bytes", BlockSize)
		return
	}
	cipher, err := aes.NewCipher(key)
	// AES 128 will only accept keys of 16 bytes (128bits) length.
	if err != nil {
		return plainText, err
	}
	plainText = make([]byte, len(cipherText))
	for bs := 0; bs < len(cipherText); bs = bs + BlockSize {
		be := bs + BlockSize
		cipher.Decrypt(plainText[bs:be], cipherText[bs:be])
	}
	return plainText, err
}

func Encrypt(plaintext, key []byte) (cipherText []byte, err error) {
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = pkcs7.Pad(plaintext, BlockSize)
	}
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return cipherText, err
	}
	cipherText = make([]byte, len(plaintext))
	for bs := 0; bs < len(plaintext); bs = bs + BlockSize {
		be := bs + BlockSize
		cipher.Encrypt(cipherText[bs:be], plaintext[bs:be])
	}
	return cipherText, err
}
