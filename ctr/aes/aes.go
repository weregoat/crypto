package aes

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"gitlab.com/weregoat/crypto/util"
	"log"
)

const blockSize = aes.BlockSize

type CTR struct {
	key [16]byte
	nonce [8]byte
}

func Decrypt(key, cipherText []byte) (string,error) {
	var iv = make([]byte, blockSize)
	var nonce = make([]byte, 8)
	copy(iv,nonce) // first 8 bytes
	var plainText []byte
	start := 0
	var count uint64 = 0
	for {
		if start >= len(cipherText) {
			break
		}
		end := start + blockSize
		if end > len(cipherText) {
			end = len(cipherText)
		}
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, count)
		if err != nil {
			log.Panic(err)
			return "", err
		}
		counter := buf.Bytes()
		for i:=0; i < len(counter); i++ {
			iv[8+i] = counter[i]
		}
		cipherBlock := cipherText[start:end]
		plainBlock := make([]byte, len(cipherBlock))
		err = processBlock(key, iv, cipherBlock, plainBlock)
		if err != nil {
			log.Panic(err)
			return "", err
		}
		plainText = append(plainText, plainBlock...)
		count++
		start = end
	}
	return string(plainText), nil
}

func Encrypt(key []byte, text string) ([]byte,error) {
	var iv = make([]byte, blockSize)
	var nonce = make([]byte, 8)
	copy(iv,nonce) // first 8 bytes
	var plainText = []byte(text)
	var cipherText []byte
	start := 0
	var count uint64 = 0
	for {
		if start >= len(plainText) {
			break
		}
		end := start + blockSize
		if end > len(plainText) {
			end = len(plainText)
		}
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, count)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		counter := buf.Bytes()
		for i:=0; i < len(counter); i++ {
			iv[8+i] = counter[i]
		}
		plainBlock := plainText[start:end]
		cipherBlock := make([]byte, len(plainBlock))
		err = processBlock(key, iv, plainBlock, cipherBlock)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		cipherText = append(cipherText, cipherBlock...)
		count++
		start = end
	}
	return cipherText, nil
}

func processBlock(key, iv, src, dst []byte) error {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	out := make([]byte, blockSize)
	cipher.Encrypt(out, iv)
	xor, err := util.FixedXORBytes(out[0:len(src)],src)
	if err != nil {
		return err
	}
	for i := 0; i < len(xor); i++ {
		dst[i] = xor[i]
	}
	return nil
}