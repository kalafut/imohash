package imohash

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/tylerb/is.v1"
)

var tempDir string

func TestMain(m *testing.M) {
	flag.Parse()

	// Make a temp area for test files
	tempDir, _ = ioutil.TempDir(os.TempDir(), "imohash_test_data")
	ret := m.Run()
	os.RemoveAll(tempDir)
	os.Exit(ret)
}

func TestCustom(t *testing.T) {
	const sampleFile = "sample"
	var hash [Size]byte
	var err error

	is := is.New(t)

	sampleSize := 3
	sampleThreshold := 45
	imo := NewCustom(sampleSize, sampleThreshold)

	// empty file
	ioutil.WriteFile(sampleFile, []byte{}, 0666)
	hash, err = imo.SumFile(sampleFile)
	is.NotErr(err)
	is.Equal(hash, [Size]byte{})

	// small file
	ioutil.WriteFile(sampleFile, []byte("hello"), 0666)
	hash, err = imo.SumFile(sampleFile)
	hashStr := fmt.Sprintf("%x", hash)
	is.Equal(hashStr, "05d8a7b341bd9b025b1e906a48ae1d19")

	/* boundary tests using the custom sample size */
	size := sampleThreshold

	// test that changing the gaps between sample zones does not affect the hash
	data := bytes.Repeat([]byte{'A'}, size)
	ioutil.WriteFile(sampleFile, data, 0666)
	h1, _ := imo.SumFile(sampleFile)

	data[sampleSize] = 'B'
	data[size-sampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h2, _ := imo.SumFile(sampleFile)
	is.Equal(h1, h2)

	// test that changing a byte on the edge (but within) a sample zone
	// does change the hash
	data = bytes.Repeat([]byte{'A'}, size)
	data[sampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h3, _ := imo.SumFile(sampleFile)
	is.NotEqual(h1, h3)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size/2] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h4, _ := imo.SumFile(sampleFile)
	is.NotEqual(h1, h4)
	is.NotEqual(h3, h4)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size/2+sampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h5, _ := imo.SumFile(sampleFile)
	is.NotEqual(h1, h5)
	is.NotEqual(h3, h5)
	is.NotEqual(h4, h5)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size-sampleSize] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h6, _ := imo.SumFile(sampleFile)
	is.NotEqual(h1, h6)
	is.NotEqual(h3, h6)
	is.NotEqual(h4, h6)
	is.NotEqual(h5, h6)

	// test that changing the size changes the hash
	data = bytes.Repeat([]byte{'A'}, size+1)
	ioutil.WriteFile(sampleFile, data, 0666)
	h7, _ := imo.SumFile(sampleFile)
	is.NotEqual(h1, h7)
	is.NotEqual(h3, h7)
	is.NotEqual(h4, h7)
	is.NotEqual(h5, h7)
	is.NotEqual(h6, h7)

	// test sampleSize < 1
	imo = NewCustom(0, size)
	data = bytes.Repeat([]byte{'A'}, size)
	ioutil.WriteFile(sampleFile, data, 0666)
	hash, _ = imo.SumFile(sampleFile)
	hashStr = fmt.Sprintf("%x", hash)
	is.Equal(hashStr, "2d9123b54d37e9b8f94ab37a7eca6f40")

	os.Remove(sampleFile)
}

// Test that the top level functions are the same as custom
// functions using the spec defaults.
func TestDefault(t *testing.T) {
	const sampleFile = "sample"
	var h1, h2 [Size]byte
	var testData []byte

	is := is.New(t)

	for _, size := range []int{100, 131071, 131072, 50000} {
		imo := NewCustom(16384, 131072)
		testData = M(size)
		is.Equal(Sum(testData), imo.Sum(testData))
		ioutil.WriteFile(sampleFile, []byte{}, 0666)
		h1, _ = SumFile(sampleFile)
		h2, _ = imo.SumFile(sampleFile)
		is.Equal(h1, h2)
	}
	os.Remove(sampleFile)
}
