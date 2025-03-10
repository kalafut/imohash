// imosum is a sample application using imohash. It will calculate and report
// file hashes in a format similar to md5sum, etc.
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/kalafut/imohash"
)

func processFile(path string) {
	hash, err := imohash.SumFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%016x  %s\n", hash, path)
}

func walkDirFunction(path string, dirEntry fs.DirEntry, err error) error {
   if err != nil {
      return err
   }
   if ! dirEntry.IsDir() {
      processFile(path)
   }
   return nil
}

func processDir(path string) {
	if err := filepath.WalkDir(path, walkDirFunction); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	paths := flag.Args()

	if len(paths) == 0 {
		fmt.Println("imosum - Prints the imohash from the file system")
		fmt.Println("USAGE: imosum path1 path2...")
		fmt.Println("If directories are provided, all files within them will be processed.")
		os.Exit(0)
	}

	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			processDir(path)
		} else {
			processFile(path)
		}
	}
}
