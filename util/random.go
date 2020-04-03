package util

import "crypto/rand"

// GenerateRandomKey returns a slice of random bytes of the given length.
func RandomBytes(size int) ([]byte, error) {
	var key []byte
	if size <= 0 {
		return key, nil
	}
	key = make([]byte, size)
	_,err := rand.Read(key)
	return key, err
}
