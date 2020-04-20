package ch16

import (
	"bytes"
	"fmt"
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"log"
)

func PoisonCipherText(o Oracle) []byte {
	blockSize := cbc.BlockSize
	p1 := bytes.Repeat([]byte{'A'}, blockSize)
	p2 := bytes.Repeat([]byte{'B'}, blockSize)
	want := []byte(";admin=true;1234")   // 1234 added to reach block size
	plaintext := string(p1) + string(p2) // Innocent enough plaintext
	c := o.Encrypt(plaintext)            // We get the innocent ciphertext
	//cipherPrint(c)
	blocks := util.Split(c, blockSize)
	c1 := blocks[2]                       // Get C1 generated with P1
	e2, err := util.FixedXORBytes(c1, p2) // Calculate E2; E2 = C1^P2
	check(err)
	inject, err := util.FixedXORBytes(e2, want) // Calculate C1'
	check(err)
	blocks[2] = inject // Replace C1 with C1' in the ciphertext
	var poisoned []byte
	for _, block := range blocks {
		poisoned = append(poisoned, block...)
	}
	//cipherPrint(poisoned)
	return poisoned // Send back the poisoned ciphertext
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cipherPrint(cipher []byte) {
	blocks := util.Split(cipher, cbc.BlockSize)
	for _, block := range blocks {
		fmt.Printf("%x\n", block)
	}
	fmt.Printf("---\n")
}

func plainPrint(plainText []byte) {
	blocks := pkcs7.Split(plainText, cbc.BlockSize)
	for _, block := range blocks {
		fmt.Printf("%+q\n", block)
	}
	fmt.Printf("---\n")
}
