package ch17

import (
	"bytes"
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
)

var blockSize = cbc.BlockSize

func getPlaintextBlock(o Oracle, c1, c2 []byte) []byte {
	var inject = make([]byte, blockSize) // The ciphertext we want to inject in place of C1 to
										 // control P2' and determine D2
	copy(inject, c1) // We can start out as C1 (it doesn't really matter)
	p2 := bytes.Repeat([]byte{'?'}, blockSize) // We store the plaintext here
	d2, _ := util.RandomBytes(blockSize) // We store our guess at the decrypted block here
	for pos := blockSize-1; pos >= 0; pos-- { // Start at the last byte of the block
		pad := blockSize-pos // The number of bytes we should pad
		padInjectBlock(inject, d2, pad) // Pad the inject block
		for j:=0; j < 256; j++ {
			inject[pos] = byte(j) // Try each possible byte value
			cipherText := append(inject, c2...) // Build the ciphertext to submit (2 blocks)
			if o.CheckPadding(cipherText) { // If the oracle said it was padded
				d2[pos] = inject[pos]^byte(pad) // The D2 byte value is D2[i] == P2'[i] ^ I[i]
				p2[pos] = d2[pos]^c1[pos] // Now that we have D2 we can XOR it to the original first block C1
										  // and extract the original plaintext P2
			}
		}
	}
	return p2
}

func padInjectBlock(inject, d2 []byte, pad int) {
	padValue := byte(pad)
	for i:=blockSize-pad; i < blockSize; i++ {
		inject[i] = d2[i]^padValue
	}
}

func attackOracle(o Oracle, cipherText, iv []byte) []byte {
	var plainText []byte

	var c1 = make([]byte, blockSize) // C1 is the block before the one we try to decipher
	var c2 = make([]byte, blockSize) // C2 is the block we try to decipher
	copy(c1, iv) // We start to try to decipher the first block of ciphertext using IV as
				// the previous block (only used for XOR).
	blocks := util.LazySplit(cipherText, blockSize)
	for i:=0; i < len(blocks); i++ {
		copy(c2, blocks[i]) // The block we want to decipher
		p2 := getPlaintextBlock(o, c1, c2) // P2 is the plaintext from block C2
		plainText = append(plainText, p2...)
		copy(c1, c2) // We use the C2 block we just deciphered as C1 for the next iteration
	}
	return pkcs7.RemovePadding(plainText)
}
