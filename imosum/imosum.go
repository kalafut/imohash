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
	//files, err := filepath.Glob(flag.Arg(0))
	//if err != nil {
	//	log.Fatal(err)
	//}

	for _, file := range files {
		hash, err := imohash.HashFilename(file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%016x  %s\n", hash, file)
	}

}
