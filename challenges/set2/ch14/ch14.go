package ch14

import (
	"bytes"
	"fmt"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/util"
	"log"
)

var BlockSize = ecb.BlockSize // We know this or we could get it

func GetSameBlockStart(cipherText []byte, blockSize int) int {
	if len(cipherText)%blockSize != 0 {
		return -1
	}
	fmt.Println(len(cipherText))
	i := 0
	for {
		if (i+2)*blockSize > len(cipherText) {
			break
		}
		b1 := cipherText[i*blockSize:(i+1)*blockSize]
		b2 := cipherText[(i+1)*blockSize:(i+2)*blockSize]
		if bytes.Equal(b1, b2) {
			return i*blockSize
		}
		i++
	}
	return -1

}

func GetPrefixLength(o Oracle) int {
	sameBlock, err := util.RandomBytes(BlockSize)
	if err != nil {
		log.Fatal(err)
	}
	var prefix []byte
	for i:=0; i < BlockSize; i++ {
		prefix = append(prefix, 'A')
		chosenText := append(prefix, bytes.Repeat(sameBlock, 2)...)
		cipherText := o.Encrypt(chosenText)
		position := GetSameBlockStart(cipherText, BlockSize)
		if position != -1 {
			return len(cipherText[0:position])-len(prefix)
		}
	}
	return -1
}