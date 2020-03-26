package util

import (
	"fmt"
	"math/bits"
)

// HammingDistance calculates the Hamming distance between two byte slices.
func HammingDistance(a, b []byte) (int, error) {

	if len(a) != len(b) {
		return 0, fmt.Errorf("supplied byte slices are of different length: %d vs %d", len(a), len(b))
	}

	/*
		This is quite simple code, I am not sure how performing it is.
		One reads a lot of alternatives with more bit manipulation and shifts,
		but, for the moment I'll stick to simple.
	*/
	var distance int
	for i := 0; i < len(a); i++ {
		xor := a[i] ^ b[i]
		distance += bits.OnesCount8(xor)
	}
	return distance, nil
}
