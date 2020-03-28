package util

import (
	"crypto/aes"
	"fmt"
)

/* https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#ECB */

func DecryptAES128ECB(cipherText, key []byte) (plainText []byte, err error) {
	blockSize := 16 // 16 bytes == 128 bits
	if len(key) != blockSize {
		err = fmt.Errorf("key size must be 16 bytes")
		return
	}
	if len(cipherText)%blockSize != 0 {
		err = fmt.Errorf("cipherText length must be a multiple of 16 bytes")
		return
	}
	cipher, err := aes.NewCipher(key)
	// AES 128 will only accept keys of 16 bytes (128bits) length.
	if err != nil {
		return plainText, err
	}
	plainText = make([]byte, len(cipherText))
	for bs := 0; bs < len(cipherText); bs = bs + blockSize {
		be := bs + blockSize
		cipher.Decrypt(plainText[bs:be], cipherText[bs:be])
	}
	return plainText, err
}
