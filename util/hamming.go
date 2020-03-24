package util

import (
	"fmt"
	"math/bits"
)

func HammingDistance(a, b []byte) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("supplied arguments are of different length: %d vs %d", len(a), len(b))
	}

	/*
	This is quite simple code, but not that performing.
	May need to rewrite it.
	 */
	distance := 0
	for i := 0; i < len(a); i++ {
			xor := a[i] ^ b[i]
			distance += bits.OnesCount8(xor)
	}
	return distance, nil
}
