package util

import (
	"bytes"
	"fmt"
	"log"
	"sort"
)

// MinBlockSize is the default minimum block size
const MinBlockSize = 2

// MaxBlockSize is the default max block size
const MaxBlockSize = 60

// NormalisedDistance stores the mean and average distances.
// It's mostly here because I wanted to experiment with means and mean and
// see if it made any difference, and because it's easier to sort.
type NormalisedDistance struct {
	BlockSize int
	Average   float64 // Normalised and Averaged distance
	Median    float64 // Use the mean instead
}

// LazySplit splits a byte slice into n byte slices of _blockSize_.
// If the source length is not a multiple of of blocksize it will leave the
// remainder out.
func LazySplit(src []byte, blockSize int) [][]byte {
	numberOfBlocks := len(src) / blockSize
	blocks := make([][]byte, numberOfBlocks)
	for i := 0; i < numberOfBlocks; i++ {
		begin := i * blockSize
		end := (i + 1) * blockSize
		if end > len(src) {
			break
		}
		blocks[i] = src[begin:end]
	}
	return blocks
}

// Split splits a byte slice into n byte slices of _blockSize_.
// It will return an error if the source length is not multiple
// of block size.
func Split(src []byte, blockSize int) ([][]byte, error) {
	var blocks [][]byte
	if len(src)%blockSize != 0 {
		err := fmt.Errorf(
			"src length %d is not a multiple of block size %d",
			len(src), blockSize,
			)
		return blocks, err
	}
	numberOfBlocks := len(src) / blockSize
	blocks = make([][]byte, numberOfBlocks)
	for i := 0; i < numberOfBlocks; i++ {
		begin := i * blockSize
		end := (i + 1) * blockSize
		if end > len(src) {
			break
		}
		blocks[i] = src[begin:end]
	}
	return blocks, nil
}

// Transpose returns a slice of slice bytes of the first byte of each block
// and the second byte of each block and so forth...
func Transpose(cipherText []byte, keySize int) [][]byte {
	var tBlocks = make([][]byte, keySize) // List of transposed blocks
	blocks := LazySplit(cipherText, keySize)
	for _, block := range blocks {
		for i := 0; i < keySize; i++ {
			if i < len(block) {
				tBlocks[i] = append(tBlocks[i], block[i])
			}
		}
	}
	return tBlocks
}

// GetBlockDistances returns the normalised distances of blocks of different sizes of the cipherText.
func GetBlockDistances(cipherText []byte, min, max int) []NormalisedDistance {
	var blockDistances []NormalisedDistance
	if min <= 0 {
		min = MinBlockSize
	}
	if max <= 0 {
		max = MaxBlockSize
	}
	for size := min; size <= max; size++ {
		// We need at least two blocks (in reality we need more, but this is the minimum).
		if len(cipherText) < size*2 {
			break
		}
		var distance = NormalisedDistance{BlockSize: size}
		// Split the cipherText in blocks of the given size
		blocks := LazySplit(cipherText, size)
		// Now we go through all the blocks and calculate the Hamming distance
		// between each block and the next
		var distances []float64 // We'll store here the normalised Hamming distances for median and mean calculation
		for i := 0; i < len(blocks)-1; i++ {
			for j := i + 1; j < len(blocks); j++ {
				distance, err := HammingDistance(blocks[i], blocks[j])
				if err != nil {
					log.Fatal(err)
				}
				distances = append(distances, float64(distance)/float64(size))
			}
		}
		// Mean of all the normalised Hamming distances for the block size:
		sum := 0.0
		for _, distance := range distances {
			sum += distance
		}
		distance.Average = sum / float64(len(distances))
		// Median
		sort.Float64s(distances)
		if len(distances)%2 == 0 {
			middle := len(distances) / 2
			distance.Median = (distances[middle+1] + distances[middle-1]) / 2
		} else {
			distance.Median = distances[len(distances)/2]
		}
		blockDistances = append(blockDistances, distance)
	}
	return blockDistances
}

// HasRepeatingBlocks returns if there is, in a supposedly ciphertext, at least
// one block that is repeated at least once (i.e. there are, at least, two
// identical blocks; could be more than one block, could be repeated more
// than once).
func HasRepeatingBlocks(src []byte, blockSize int) bool {
	blocks := LazySplit(src, blockSize)
	for i := 0; i < len(blocks)-1; i++ {
		for j := i + 1; j < len(blocks); j++ {
			if bytes.Equal(blocks[i], blocks[j]) {
				return true
			}
		}
	}
	return false
}
