package pkcs7

import (
	"bytes"
)

func Pad(src []byte, size int) []byte {
	padded := src
	if size <= 0 {
		return src
	}
	rest := len(src) % size
	paddingLength := size - rest
	/* Adds a full padded block if no padding is required */
	if paddingLength == 0 {
		paddingLength = size
	}
	value := byte(paddingLength)
	padding := bytes.Repeat([]byte{value}, paddingLength)
	return append(padded, padding...)
}

func RemovePadding(src []byte) []byte {
	// By default the last byte position is the end of the slice
	last := len(src) - 1
	// A slice smaller than 1 makes little sense
	if last < 0 {
		return src
	}
	// Pick the last byte value, it should the the number of bytes in the padding
	count := int(src[last])
	// It cannot be padded if the slice length is less than the padding length
	if len(src) < count {
		return src
	}
	// If the text is padded is should start at len - endByte
	padStart := len(src) - count

	// Pick count bytes from the end and they all should have count value
	for i := padStart; i < len(src); i++ {
		if int(src[i]) != count {
			return src
		}
	}
	return src[:padStart]
}

// Split splits a byte slice into n byte slices of _blockSize_.
func Split(src []byte, blockSize int) [][]byte {
	var blocks [][]byte
	i := 0
	for {
		begin := i * blockSize
		if begin >= len(src) {
			break
		}
		end := (i + 1) * blockSize
		if end > len(src) {
			end = len(src)
		}
		// Padding
		block := src[begin:end]
		if len(block)%blockSize != 0 {
			block = Pad(block, blockSize)
		}
		blocks = append(blocks, block)
		i++
	}
	return blocks
}