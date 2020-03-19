package util

import "bytes"

func Match(plainText []byte, punctuation []byte, dictionary [][]byte) int {
	if len(plainText) == 0 {
		return 0
	}
	/* always use lowercase for simplicity */
	lPlainText := bytes.ToLower(plainText)
	var words [][]byte
	if len(punctuation) == 0 {
		words = dictionary
	} else {
		for _,word := range dictionary {
			for _,p := range punctuation {
				word = append(word, p)
				words = append(words, word)
			}
		}
	}
	total := 0
	for _, word := range words {
		if len(word) > 0 {
			lWord := bytes.ToLower(word)
			total += bytes.Count(lPlainText, lWord)
		}
	}
	return total
}
