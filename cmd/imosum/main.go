// imosum is a sample application using imohash. It will calculate and report
// file hashes in a format similar to md5sum, etc.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kalafut/imohash"
)

func processFile(path string) {
	hash, err := imohash.SumFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%016x  %s\n", hash, path)
}

func processWalkFunc(path string, dirEntry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !dirEntry.IsDir() {
		processFile(path)
	}
	return nil
}

func checkFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	malformedLines := 0
	failedHashes := 0

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "  ", 2)

		if len(parts) != 2 {
			malformedLines += 1
			continue
		}

		readableResult := "FAILED"
		hash, err := imohash.SumFile(parts[1])
		if err != nil {
			failedHashes += 1
		} else {
			hashStr := fmt.Sprintf("%016x", hash)

			if hashStr == parts[0] {
				readableResult = "OK"
			} else {
				failedHashes += 1
			}
		}

		fmt.Printf("%s: %s\n", parts[1], readableResult)
	}

	if failedHashes > 0 {
		checksumStr := "checksum"
		if failedHashes > 1 {
			checksumStr += "s"
		}
		fmt.Fprintf(os.Stderr, "imosum: WARNING: %d computed %s did NOT match\n", failedHashes, checksumStr)
	}

	if malformedLines > 0 {
		lineStr := "line is"
		if malformedLines > 1 {
			lineStr = "lines are"
		}
		fmt.Fprintf(os.Stderr, "imosum: WARNING: %d %s improperly formatted\n", malformedLines, lineStr)
	}
}

func checkWalkFunc(path string, dirEntry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !dirEntry.IsDir() {
		checkFile(path)
	}
	return nil
}

func walkDir(path string, fn fs.WalkDirFunc) {
	if err := filepath.WalkDir(path, fn); err != nil {
		log.Fatal(err)
	}
}

func main() {
	chkFlag := flag.Bool("c", false, "read hashes from the FILEs and check them")
	flag.Parse()
	paths := flag.Args()

	if len(paths) == 0 {
		fmt.Println("imosum - Prints the imohash from the file system")
		fmt.Println("USAGE: imosum [-c] path1 path2...")
		fmt.Println("If directories are provided, all files within them will be processed.")
		fmt.Println("Pass -c to read hashes from the FILEs and check them.")
		os.Exit(0)
	}

	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			if *chkFlag {
				walkDir(path, checkWalkFunc)
			} else {
				walkDir(path, processWalkFunc)
			}
		} else {
			if *chkFlag {
				checkFile(path)
			} else {
				processFile(path)
			}
		}
	}
}
