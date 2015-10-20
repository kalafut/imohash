package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kalafut/imohash"
)

func main() {
	flag.Parse()
	files := flag.Args()

	for _, file := range files {
		hash, err := imohash.SumFile(file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%016x  %s\n", hash, file)
	}

}
