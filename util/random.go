package util

import (
	rand "crypto/rand"
	"log"
	"math/big"
)

// GenerateRandomKey returns a slice of random bytes of the given length.
func RandomBytes(size int) ([]byte, error) {
	var key []byte
	if size <= 0 {
		return key, nil
	}
	key = make([]byte, size)
	_, err := rand.Read(key)
	return key, err
}

func RandomInt(min, max int) int {
	a := int64(min)
	b := int64(max)
	diff := big.NewInt(b - a)
	r, err := rand.Int(rand.Reader, diff)
	if err != nil {
		// The RNG must have failed
		log.Fatal(err)
	}
	c := int(r.Int64()) + min
	return c
}

// RandomPad returns a byte slice padded at both ends with a random number of
// bytes picked between min and max each time.
func RandomPad(src []byte, min, max int) []byte {
	before, err := RandomBytes(RandomInt(min, max))
	if err != nil {
		log.Fatal(err)
	}
	after, err := RandomBytes(RandomInt(min, max))
	if err != nil {
		log.Fatal(err)
	}
	src = append(before, src...)
	src = append(src, after...)
	return src
}
