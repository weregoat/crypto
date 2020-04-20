package ch14

import (
	"bytes"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"log"
)

var blockSize = ecb.BlockSize // We know this or we could get it
var fillUp = []byte{'A'}      // We use this for chosen text fill-ups

func GetSameBlockStart(src []byte, blockSize int) int {
	cipherText := make([]byte, len(src))
	copy(cipherText, src)
	if len(cipherText)%blockSize != 0 {
		log.Fatalf("ciphertext length %d not multiple of blocksize %d", len(cipherText), blockSize)
	}
	blocks := util.LazySplit(cipherText, blockSize)
	for i := 0; i < len(blocks)-1; i++ {
		if bytes.Equal(blocks[i], blocks[i+1]) {
			return i * blockSize
		}
	}
	return -1

}

func GetPrefixLength(o Oracle) int {
	sameBlock := bytes.Repeat([]byte{'A'}, blockSize)
	for i := 0; i < blockSize; i++ {
		prefix := bytes.Repeat([]byte{'B'}, i)
		chosenText := append(prefix, bytes.Repeat(sameBlock, 2)...)
		cipherText := o.Encrypt(chosenText)
		position := GetSameBlockStart(cipherText, blockSize)
		if position != -1 {
			if verifyPrefix(o, i) == position {
				return position - len(prefix)
			}
		}
	}
	return -1
}

func LookupTable(oracle Oracle, knownPlaintext []byte, begin int) map[string]byte {
	var table = make(map[string]byte) // using a string for key will make lookup simpler
	for i := 0; i <= 255; i++ {
		p := append(knownPlaintext, byte(i))
		cipherText := oracle.Encrypt(p)
		table[string(cipherText[begin:begin+blockSize])] = byte(i)
		//fmt.Printf("%x=>%x\n", cipherText[skipTo:skipTo+BlockSize], byte(i))
	}
	return table
}

func CPA(oracle Oracle) []byte {
	// The slice where we store the plaintext
	var plainText []byte
	// How many bytes in the prefix
	prefixLength := GetPrefixLength(oracle)
	/*
		16byte plaintext => 16bytes ciphertext, so we can calculate how to
		restrict the prefix + some chosen text into a few block of ciphertext
		we can then ignore and solve this as #12.
	*/
	// How many bytes of padding should we add so that the prefix + padding
	// results in the same blocks of ciphertext (so we can ignore them?)
	prefixPaddingLength := blockSize - prefixLength%blockSize
	// We can ignore the ciphertext before this
	plainTextBegin := prefixLength + prefixPaddingLength
	// Creates a byte slice to use for padding the prefix
	prefixPadding := bytes.Repeat(fillUp, prefixPaddingLength)
	/* I want to do things a bit differently than in #12, and try out not to use util.Split */
	// Maximum length the plaintext could be (can be shorter because of padding)
	maxPlaintext := len(oracle.Encrypt(prefixPadding)) - plainTextBegin // Base ciphertext length
	// One plaintext block at a time
	for i := 0; i < maxPlaintext; i = i + blockSize {
		// One byte at a time
		for j := 1; j <= blockSize; j++ {
			// The part of the plaintext we are using for the table
			var knownPlainText = make([]byte, blockSize-1)
			// If we don't have blocksize worth of plaintext, we add what
			// we have to the chosentext because that's what the oracle will
			// encrypt (we haven't shifted a whole block yet)
			if i < blockSize {
				knownPlainText = append(bytes.Repeat(fillUp, blockSize-j), plainText...)
			} else {
				// Otherwise just pick the last 15 bytes of the plaintext
				// Remember this is for looking up in the table.
				copy(knownPlainText, plainText[len(plainText)-(blockSize-1):])
			}
			knownPlainText = append(prefixPadding, knownPlainText...)
			// Now we build a table of the ciphertexts from the known plaintext + every byte
			table := LookupTable(oracle, knownPlainText, plainTextBegin)
			// What we are going to submit to have the plaintext shift
			chosenPlainText := append(prefixPadding, bytes.Repeat(fillUp, blockSize-j)...) // first one short of blocksize, then two... to have the plaintext shift left
			// We submit the chosen plaintext
			cipherText := oracle.Encrypt(chosenPlainText)
			// We need to select the right block of the ciphertext
			blockStart := plainTextBegin + i
			cipherBlock := cipherText[blockStart : blockStart+blockSize]
			// Lookup the cipherblock in the table
			plainTextByte, ok := table[string(cipherBlock)]
			// If we could find the ciphertext in the table
			if ok {
				// Add the byte to the plaintext
				plainText = append(plainText, plainTextByte)
			}
			//fmt.Printf("%+q\n", plainText)
			// If we don't find it, it could be because of padding
			// Not covering that at the moment, will fix.
		}
	}
	return pkcs7.RemovePadding(plainText)
}

func cipherPrint(cipher []byte) {
	blocks, err := util.Split(cipher, blockSize)
	if err != nil {
		log.Println(err)
	}
	for _, block := range blocks {
		log.Printf("%x", block)
	}
	log.Printf("---")
}

func verifyPrefix(o Oracle, i int) int {
	sameBlock := bytes.Repeat([]byte{'C'}, blockSize)
	prefix := bytes.Repeat([]byte{'D'}, i)
	chosenText := append(prefix, bytes.Repeat(sameBlock, 2)...)
	cipherText := o.Encrypt(chosenText)
	return GetSameBlockStart(cipherText, blockSize)
}
