package main

import (
	"encoding/base64"
	"gitlab.com/weregoat/crypto/util"
	"io/ioutil"
	"log"
	"sort"
)

type BlockDistance struct {
	Size                      int
	AverageNormalisedDistance float64 // Normalised and Averaged distance
	MedianNormalisedDistance  float64 // Use the mean instead
}

func main() {
	encoded, err := ioutil.ReadFile("6.txt")
	if err != nil {
		log.Fatal(err)
	}
	cypherText, _ := base64.StdEncoding.DecodeString(string(encoded))
	distances := GetBlockDistances(cypherText, 40)
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].MedianNormalisedDistance < distances[j].MedianNormalisedDistance
	})



}

/*
For each KEYSIZE, take the first KEYSIZE worth of bytes, and the second KEYSIZE worth of bytes, and find the edit distance between them. Normalize this result by dividing by KEYSIZE.
The KEYSIZE with the smallest normalized edit distance is probably the key. You could proceed perhaps with the smallest 2-3 KEYSIZE values. Or take 4 KEYSIZE blocks instead of 2 and average the distances.

Note: the above instructions were confusing to me, because it seems that it should be enough to use the distance of the
first two blocks. It is, actually, not so.
https://en.wikipedia.org/wiki/Vigenère_cipher
Possible explanation on the reason behind the suggested method:
https://crypto.stackexchange.com/questions/8115/repeating-key-xor-and-hamming-distance
https://laconicwolf.com/2018/06/30/cryptopals-challenge-6-break-repeating-key-xor/

I don't know if this was intentional, or it's just me having problem reading their instructions.
Anyway, it made me read a bit on how you are supposed to do it, and I came out with this:
*/
func GetBlockDistances(cypherText []byte, maxKeySize int) []BlockDistance {
	var distances []BlockDistance
	for keySize := 2; keySize <= maxKeySize; keySize++ {
		// We need at least two blocks (in reality we need more, but this is the minimum).
		if len(cypherText) < maxKeySize*2 {
			break
		}
		var distance = BlockDistance{Size: keySize}
		// Split the cyphertext in blocks of the given size
		blocks := Split(cypherText, keySize)
		// Now we go through all the blocks and calculate the Hamming distance
		// between each block and the next
		var nhd []float64 // We'll store here the normalised Hamming distances for median and mean calculation
		for i := 0; i < len(blocks) - 1 ; i++ {
			for j:=i+1; j < len(blocks); j++ {
				hd, err := util.HammingDistance(blocks[i], blocks[j])
				if err != nil {
					log.Fatal(err)
				}
				nhd = append(nhd, float64(hd)/float64(keySize))
			}
		}
		// Mean Hamming distance:
		sum := 0.0
		for _,hd := range nhd {
			sum += hd
		}
		distance.AverageNormalisedDistance = sum/float64(len(nhd))
		// Median
		sort.Float64s(nhd)
		if len(nhd)%2 == 0 {
			middle := len(nhd)/2
			distance.MedianNormalisedDistance = (nhd[middle+1] + nhd[middle-1])/2
		} else {
			distance.MedianNormalisedDistance = nhd[len(nhd)/2]
		}
		distances = append(distances, distance)
	}
	return distances
}

/*
Now that you probably know the KEYSIZE: break the ciphertext into blocks of KEYSIZE length.
Now transpose the blocks: make a block that is the first byte of every block, and a block that is the second byte of every block, and so on.
Solve each block as if it was single-character XOR. You already have code to do this.
For each block, the single-byte XOR key that produces the best looking histogram is the repeating-key XOR key byte for that block. Put them together and you have the key.
*/

func Split(src []byte, blockSize int) [][]byte {
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
