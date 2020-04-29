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
	return plainText, nil
}

func Encrypt(src, key []byte) (cipherText []byte, err error) {
	var plaintext []byte
	// Only pad the plaintext if size not a multiple of the blocksize
	if len(src)%BlockSize != 0 {
		plaintext = pkcs7.Pad(src, BlockSize)
	} else {
		plaintext = make([]byte, len(src))
		copy(plaintext, src)
	}
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return cipherText, err
	}
	cipherText = make([]byte, len(plaintext))
	i := 0
	for {
		bs := i * BlockSize
		be := (i + 1) * BlockSize
		if be > len(plaintext) {
			break
		}
		cipher.Encrypt(cipherText[bs:be], plaintext[bs:be])
		i++
	}
	return cipherText, err
}
