package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"gitlab.com/weregoat/crypto/util"
	"log"
	"os"
)

/*
Write a function that takes two equal-length buffers and produces their XOR combination.

If your function works properly, then when you feed it the string:

1c0111001f010100061a024b53535009181c
... after hex decoding, and when XOR'd against:

686974207468652062756c6c277320657965
... should produce:

746865206b696420646f6e277420706c6179
 */

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		log.Fatal("hex strings A and B are required")
	}
	a, err := hex.DecodeString(args[1])
	if err != nil {
		log.Fatal(err)
	}
	b, err := hex.DecodeString(args[2])
	if err != nil {
		log.Fatal(err)
	}
	c, err := util.FixedXORBytes(a, b)
	if err != nil {
		log.Fatal(c)
	}
	fmt.Printf("A String: %s '%s'\n", hex.EncodeToString(a), a)
	fmt.Printf("B String: %s '%s'\n", hex.EncodeToString(b), b)
	fmt.Printf("Result:   %s '%s'\n", hex.EncodeToString(c), c)
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s [hex string] [hex string]\n", os.Args[0])
	flag.PrintDefaults()
}
