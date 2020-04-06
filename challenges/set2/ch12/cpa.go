package ch12

import (
	"bytes"
	"gitlab.com/weregoat/crypto/util"
)

// GuessBlockSize tries to guess the cipher blocksize by adding a byte at a time
// to the plaintext until the oracle returns a ciphertext with an extra block.
func GuessBlockSize(oracle Oracle) int {
	plainText := make([]byte, 0)
	// This is the original length of the ciphertext
	orig := len(oracle.Encrypt(plainText))
	for {
		// Add one byte to the chosen plaintext
		plainText = append(plainText, 'A')
		cipherText := oracle.Encrypt(plainText)
		// If the length of the cipherText has increased is
		// because a new block has been added.
		if len(cipherText) > orig {
			// The block size is the difference in length
			return len(cipherText) - orig
		}

	}
}

// IsECB tries to verify that the encryption mode is ECB by detecting block
// repetitions.
func IsECB(oracle Oracle, blockSize int) bool {
	chosenPlaintext := bytes.Repeat([]byte{'A'}, blockSize*4) // We determined earlier than 4 is a good enough number
	cipherText := oracle.Encrypt(chosenPlaintext)
	return util.HasRepeatingBlocks(cipherText, blockSize)
}

func LookupTable(oracle Oracle, knownPlaintext []byte, begin, end int) map[string]byte {
	var table = make(map[string]byte) // using a string for key will make lookup simpler
	for i:=0; i <= 255; i++ {
		chosenPlaintext := knownPlaintext[:15]
		chosenPlaintext = append(chosenPlaintext, byte(i))
		cipherText := oracle.Encrypt(chosenPlaintext)
		table[string(cipherText[begin:end])] = byte(i)
	}
	return table
}