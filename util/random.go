package util

import (
	cryptorand "crypto/rand"
	"log"
	"math/rand"
)

// GenerateRandomKey returns a slice of random bytes of the given length.
func RandomBytes(size int) ([]byte, error) {
	var key []byte
	if size <= 0 {
		return key, nil
	}
	key = make([]byte, size)
	_,err := cryptorand.Read(key)
	return key, err
}

// Not cryptographically secure
func RandomPad(src []byte, min,max int) []byte {
	before, err := RandomBytes(rand.Intn(max-min)+min)
	if err != nil {
		log.Fatal(err)
	}
	after, err := RandomBytes(rand.Intn(max-min)+min)
	if err != nil {
		log.Fatal(err)
	}
	src = append(before, src...)
	src = append(src, after...)
	return src
}