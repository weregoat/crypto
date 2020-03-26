package util

import (
	"fmt"
)

/*
Fixed XOR
Write a function that takes two equal-length buffers and produces their XOR combination.

If your function works properly, then when you feed it the string:

1c0111001f010100061a024b53535009181c
... after hex decoding, and when XOR'd against:

686974207468652062756c6c277320657965
... should produce:

746865206b696420646f6e277420706c6179
*/

func FixedXORBytes(a, b []byte) ([]byte, error) {
	n := len(a)
	if len(b) != n {
		err := fmt.Errorf("different buffer length %d<>%d", len(a), len(b))
		return nil, err
	}
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst, nil
}

func RepeatingXORBytes(plainText []byte, key []byte) []byte {
	cypherText := make([]byte, len(plainText))
	for i, j := range plainText {
		cypherText[i] = j ^ key[i%len(key)]
	}
	return cypherText
}
