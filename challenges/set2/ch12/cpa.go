package ch12

import (
	"bytes"
	"gitlab.com/weregoat/crypto/pkcs7"
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

func LookupTable(oracle Oracle, knownPlaintext []byte, blockSize int) map[string]byte {
	var table = make(map[string]byte) // using a string for key will make lookup simpler
	for i := 0; i <= 255; i++ {
		p := append(knownPlaintext, byte(i))
		cipherText := oracle.Encrypt(p)
		table[string(cipherText[0:blockSize])] = byte(i)
	}
	return table
}

func CPA(oracle Oracle) []byte {
	var plainText []byte
	blockSize := GuessBlockSize(oracle)
	shortBlock := blockSize-1 // A whole block minus one byte
	chosenByte := []byte{'A'} // It doesn't really matter
	// How many blocks the plaintext has
	cipherBlocks := len(util.Split(oracle.Encrypt([]byte{}), blockSize))
	// For each of the blocks
	for i := 0; i < cipherBlocks; i++ {
		// One byte at a time
		for j:=1; j <= blockSize; j++ {
			// The part of the plaintext we are using for the table
			var knownPlainText = make([]byte, shortBlock)
			// If we don't have blocksize worth of plaintext, we add what
			// we have to the chosentext because that's what the oracle will
			// encrypt (we haven't shifted a whole block yet)
			if len(plainText) < blockSize {
				knownPlainText = append(bytes.Repeat(chosenByte, blockSize-j), plainText...)
			} else {
				// Otherwise just pick the last blocksize-1 bytes of the plaintext
				// Remember this is for looking up in the table.
				copy(knownPlainText,plainText[len(plainText)-shortBlock:])
			}
			// Now we build a table of the ciphertexts from the known plaintext + every byte
			table := LookupTable(oracle, knownPlainText, blockSize)

			// What we are going to submit to have the plaintext shift
			chosenPlainText := bytes.Repeat(chosenByte, blockSize-j) // first one short of blocksize, then two... to have the plaintext shift left
			// We submit the chosen plaintext
			cipherText := oracle.Encrypt(chosenPlainText)
			// We need to select the right block of the ciphertext
			cipherBlock := cipherText[i*blockSize:(i+1)*blockSize]
			// Lookup the cipherblock in the table
			plainTextByte, ok := table[string(cipherBlock)]
			// If we could find the ciphertext in the table
			if ok {
				// Add the byte to the plaintext
				plainText = append(plainText, plainTextByte)
			}
			// If we don't find it, it could be because of padding
			// Not covering that at the moment, will fix.
		}
	}
	plainText, _ = pkcs7.RemovePadding(plainText)
	return plainText
}
