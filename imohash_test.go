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

func TestDefault(t *testing.T) {
	const sampleFile = "sample"
	var hash [Size]byte
	var err error

	is := is.New(t)

	// empty file
	ioutil.WriteFile(sampleFile, []byte{}, 0666)
	hash, err = SumFile(sampleFile)
	is.NotErr(err)
	is.Equal(hash, [Size]byte{})

	// small file
	hash = Sum([]byte("hello"))
	hashStr := fmt.Sprintf("%x", hash)
	is.Equal(hashStr, "05d8a7b341bd9b025b1e906a48ae1d19")

	/* boundary tests using the default sample size */
	size := SampleThreshhold

	// test that changing the gaps between sample zones does not affect the hash
	data := bytes.Repeat([]byte{'A'}, size)
	ioutil.WriteFile(sampleFile, data, 0666)
	h1, _ := SumFile(sampleFile)

	data[SampleSize] = 'B'
	data[size-SampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h2, _ := SumFile(sampleFile)
	is.Equal(h1, h2)

	// test that changing a byte on the edge (but within) a sample zone
	// does change the hash
	data = bytes.Repeat([]byte{'A'}, size)
	data[SampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h3, _ := SumFile(sampleFile)
	is.NotEqual(h1, h3)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size/2] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h4, _ := SumFile(sampleFile)
	is.NotEqual(h1, h4)
	is.NotEqual(h3, h4)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size/2+SampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h5, _ := SumFile(sampleFile)
	is.NotEqual(h1, h5)
	is.NotEqual(h3, h5)
	is.NotEqual(h4, h5)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size-SampleSize] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h6, _ := SumFile(sampleFile)
	is.NotEqual(h1, h6)
	is.NotEqual(h3, h6)
	is.NotEqual(h4, h6)
	is.NotEqual(h5, h6)

	// test that changing the size changes the hash
	data = bytes.Repeat([]byte{'A'}, size+1)
	ioutil.WriteFile(sampleFile, data, 0666)
	h7, _ := SumFile(sampleFile)
	is.NotEqual(h1, h7)
	is.NotEqual(h3, h7)
	is.NotEqual(h4, h7)
	is.NotEqual(h5, h7)
	is.NotEqual(h6, h7)

	os.Remove(sampleFile)
}

// Test the basic hash.Hash functions
func TestHashInterface(t *testing.T) {
	const sampleFile = "sample"

	is := is.New(t)

	// Test Write() and Sum()
	defaultHasher.Reset()
	defaultHasher.Write([]byte("hello"))
	base := []byte{0x55, 0x22, 0xee}
	hashStr := fmt.Sprintf("%x", defaultHasher.Sum(base))
	is.Equal(hashStr, "5522ee05d8a7b341bd9b025b1e906a48ae1d19") // matches reference murmur3 hash

	// Test that calling Sum() previously didn't affect state
	base = []byte{0x55, 0x22, 0xee}
	hashStr = fmt.Sprintf("%x", defaultHasher.Sum(base))
	is.Equal(hashStr, "5522ee05d8a7b341bd9b025b1e906a48ae1d19")

	// Test adding more data with Write()
	defaultHasher.Write([]byte(", world"))
	base = []byte{0x99}
	hashStr = fmt.Sprintf("%x", defaultHasher.Sum(base))
	is.Equal(hashStr, "990c2fac623a5ebc8e4cdcbc079642414d") // matches reference murmur3 hash

	// Test Reset()
	defaultHasher.Reset()
	base = []byte{}
	hashStr = fmt.Sprintf("%x", defaultHasher.Sum(base))
	is.Equal(hashStr, "00000000000000000000000000000000")

	// Test BlockSize() and Size()
	is.Equal(defaultHasher.BlockSize(), 1)
	is.Equal(defaultHasher.Size(), 16)
}
