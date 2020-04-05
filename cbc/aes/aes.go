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
	var cipherText []byte
	// We need to make a copy of the IV because otherwise we'll change it
	IV := make([]byte, len(iv))
	copy(IV, iv)
	// Pad, if needed, and split the plaintext into blocks
	blocks := pkcs7.Split(plainText, BlockSize)
	// Initialise the AES cipher block
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return cipherText, err
	}
	// Encrypt each block
	// By using copies we make sure we don't mess with original slices
	// It's, probably, slower than manipulating the slices directly, but
	// I find it much clearer to show how CBC works.
	for _, block := range blocks {
		// Copy each block so we can manipulate it
		src := make([]byte, BlockSize)
		copy(src, block)
		// XOR the block with the IV
		xor(src, IV)
		if err != nil {
			return cipherText, err
		}
		// Create an empty block for the plaintext
		dst := make([]byte, BlockSize)
		// Encrypt the XOR block into the plaintext block
		cipher.Encrypt(dst, src)
		cipherText = append(cipherText, dst...)
		copy(IV, dst)
	}
	return cipherText, nil
}

// Decrypt returns a plaintext from a ciphertext encrypted with AES in CBC mode.
func Decrypt(cipherText, key, iv []byte) ([]byte, error) {
	var plainText []byte
	// Make a copy of the IV to avoid changing it
	IV := make([]byte, len(iv))
	copy(IV, iv)
	// The cipherText length should be a multiple of AES blocksize (128 bits/16 bytes)
	if len(cipherText)%BlockSize != 0 {
		err := fmt.Errorf("ciphertext byte length(%d) must be a multiple of %d bytes", len(cipherText), BlockSize)
		return plainText, err
	}
	// Initialise the AES cipher
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return plainText, err
	}
	// We should not need padding.
	// PKCS7 split function would work as well, but here
	// using this for
	blocks := util.Split(cipherText, BlockSize)
	// Go through each block and decipher it.
	for _, block := range blocks {
		// Copy the block, so we don't change it
		src := make([]byte, BlockSize)
		copy(src, block)
		// Create a block to store the deciphered block into
		dst := make([]byte, BlockSize)
		// Decrypt the block of ciphertext
		cipher.Decrypt(dst, src)
		// XOR the decrypted block with the IV
		xor(dst, IV)
		if err != nil {
			return plainText, err
		}
		// Add the resulting block to the plaintext
		plainText = append(plainText, dst...)
		// Update the IV with the original cipher block
		copy(IV, src)
	}
	// Remove the padding before returning
	return pkcs7.RemovePadding(plainText), nil
}

// XOR function to be used internally
func xor(dst, src []byte) {
	for i := 0; i < len(dst); i++ {
		dst[i] = dst[i] ^ src[i]
	}
}
