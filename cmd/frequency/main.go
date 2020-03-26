package main

import (
	"bytes"
	"flag"
	"fmt"
	"gitlab.com/weregoat/crypto/util"
	"io/ioutil"
	"log"
	"sort"
)

func main() {
	letters := flag.String("letters", "", "String with all the letters to analyse")
	corpus := flag.String("corpus", "", "Text file to use for analysis")
	lowercase := flag.Bool("lowercase", true, "Lowecase the corpus before processing")
	flag.Parse()
	source, err := ioutil.ReadFile(*corpus)
	if err != nil {
		log.Fatal(err)
	}
	if *lowercase {
		source = bytes.ToLower(source)
	}

	frequencies := util.FrequencyCalculator(source, []byte(*letters))

	sort.Slice(frequencies, func(i int, j int) bool {
		return frequencies[i].Value > frequencies[j].Value
	})
	for _, frequency := range frequencies {
		fmt.Printf("{Character:%q,Value:%.5f},\n", frequency.Character, frequency.Value)
	}
}
