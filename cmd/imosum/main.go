// imosum is a sample application using imohash. It will calculate and report
// file hashes in a format similar to md5sum, etc.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kalafut/imohash"
)

func main() {
	flag.Parse()
	files := flag.Args()

	if len(files) == 0 {
		fmt.Println("imosum filenames")
		os.Exit(0)
	}

	for _, file := range files {
		hash, err := imohash.SumFile(file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%016x  %s\n", hash, file)
	}
}
