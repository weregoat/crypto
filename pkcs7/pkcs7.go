package pkcs7

import (
	"bytes"
)


/*
RFC 2315
[...]
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
[...]
 */

// Pad returns a slice of bytes from src with PKCS#7 padding at the end.
func Pad(src []byte, k int) []byte {
	// Always copy slice so we always return a copy regardless of padding or
	dst := make([]byte, len(src))
	copy(dst, src)
	// Invalid block size
	if k <= 0 {
		return dst
	}
	l := len(dst)
	n := k - (l % k)
	value := byte(n)
	padding := bytes.Repeat([]byte{value}, n)
	return append(dst, padding...)
}

// RemovePadding removes the PKCS#7 padding from a slice of bytes.
func RemovePadding(src []byte) []byte {
	// Something quite important, I think.
	// I should **always** return a copy of the slice
	// because that's what we do when we remove the padding.
	var dst = make([]byte, len(src))
	copy(dst, src)
	start := paddingStart(dst)
	if start >= 0 {
		return dst[:start]
	}
	return dst
}

// paddingStart returns the postion where the padding start. -1 if there is no
// PKCS#7 padding.
func paddingStart(src []byte) int {

	// Because even if the original plaintext was empty, it should
	// have, according to PKCS#7 rules a block of padding attached.
	if len(src) == 0 {
		return -1
	}

	// index of the last byte
	i := len(src) - 1
	// The value as integer (is really an uint, but int is easier to handle around)
	n := int(src[i])
	// The value is supposed to be the number of bytes of padding we should
	// remove. It cannot be higher than the size of the slice (although it can
	// be the same, if the plaintext is an empty string).
	if n > len(src) {
		return -1
	}
	// If the padding is correct, it should start here
	start := len(src) - n
	// Check every that each of the n-1 bytes in the slice has the same
	// value as n (we already know the last one has n)
	for j:=1; j < n; j++ {
		b := int(src[i-j])
		if b != n {
			return -1
		}
	}
	return start
}

// IsPadded returns true if the plaintext is correctly padded according to
// PKCS#7.
func IsPadded(plainText []byte) bool {
	// I don't want to go all ontological on this, but it seems to me that
	// a padded string **needs** to have a point where the padding starts
	// otherwise is not padded.
	if paddingStart(plainText) > 0 {
		return true
	}
	return false
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
