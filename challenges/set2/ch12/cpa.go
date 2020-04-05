package ch12

/*
func main() {
	secret := "Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK"
}
*/

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
