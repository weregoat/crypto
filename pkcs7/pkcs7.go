package pkcs7

import (
	"bytes"
	"fmt"
)

type NotPKCSError struct {

}

/*
   2.   Some content-encryption algorithms assume the
        input length is a multiple of k octets, where k > 1, and
        let the application define a method for handling inputs
        whose lengths are not a multiple of k octets. For such
        algorithms, the method shall be to pad the input at the
        trailing end with k - (l mod k) octets all having value k -
        (l mod k), where l is the length of the input. In other
        words, the input is padded at the trailing end with one of
        the following strings:

                 01 -- if l mod k = k-1
                02 02 -- if l mod k = k-2
                            .
                            .
                            .
              k k ... k k -- if l mod k = 0

        The padding can be removed unambiguously since all input is
        padded and no padding string is a suffix of another. This
        padding method is well-defined if and only if k < 256;
        methods for larger k are an open issue for further study.
 */

// Pad returns a slice of bytes from src with PKCS#7 padding at the end.
func Pad(src []byte, size int) []byte {
	padded := src
	// No padding
	if size <= 0 {
		return src
	}
	rest := len(padded) % size
	paddingLength := size - rest
	value := byte(paddingLength)
	padding := bytes.Repeat([]byte{value}, paddingLength)
	return append(padded, padding...)
}

// RemovePadding removes the PKCS#7 padding from a slice of bytes.
func RemovePadding(src []byte) ([]byte, error) {
	// By default the last byte position is the end of the slice
	last := len(src) - 1
	// A slice smaller than 1 makes little sense
	if last < 0 {
		return src, nil
	}
	// Pick the last byte value, it should the the number of bytes of the padding
	count := int(src[last])
	// It cannot be padded if the slice length is less than the padding length
	if len(src) < count {
		err := fmt.Errorf("slice size %d is less than padding value %d", len(src), count)
		return src, err
	}
	// If the text is padded is should start at len - endByte
	padStart := len(src) - count

	// Pick count bytes from the end and they all should have count value
	for i := padStart; i < len(src); i++ {
		if int(src[i]) != count {
			err := fmt.Errorf(
				"invalid byte value position %d; expecting %d, but found %d",
				i, count, int(src[i]),
				)
			return src, err
		}
	}
	return src[:padStart], nil
}

// Split splits a byte slice into n byte slices of _blockSize_.
func Split(src []byte, blockSize int) [][]byte {
	if len(src)%blockSize != 0 {
		src = Pad(src, blockSize)
	}
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
		blocks = append(blocks, src[begin:end])
		i++
	}
	return blocks
}
