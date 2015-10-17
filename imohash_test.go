package imohash

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/spaolacci/murmur3"
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

func TestFoldInt(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		in  int64
		out uint32
	}{
		{0x00, 0x00}, {0x01, 0x01}, {0xff, 0xff}, {0xffff, 0xffff},
		{0xffffffff, 0xffffffff}, {0x01ffffffff, 0xfeffffff}, {0x123456789abcdef0, 0xe2eaeae2},
	}

	for _, test := range tests {
		is.Equal(foldInt(test.in), test.out)
	}
}

func TestDefault(t *testing.T) {
	const sampleFile = "sample"
	var hash uint64
	var err error

	is := is.New(t)

	// empty file
	ioutil.WriteFile(sampleFile, []byte{}, 0666)
	hash, err = HashFilename(sampleFile)
	is.NotErr(err)
	is.Equal(hash, uint64(0)<<32|uint64(murmur3.Sum32([]byte{})))

	// small file
	ioutil.WriteFile(sampleFile, []byte("hello"), 0666)
	hash, err = HashFilename(sampleFile)
	is.Equal(hash, uint64(5)<<32|uint64(murmur3.Sum32([]byte{'h', 'e', 'l', 'l', 'o'})))

	/* boundary tests using the default sample size */
	size := 12290

	// test that changing the gaps between sample zones does not affect the hash
	data := bytes.Repeat([]byte{'A'}, size)
	ioutil.WriteFile(sampleFile, data, 0666)
	h1, _ := HashFilename(sampleFile)

	data[defaultSampleSize] = 'B'
	data[size-defaultSampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h2, _ := HashFilename(sampleFile)
	is.Equal(h1, h2)

	// test that changing a byte on the edge (but within) a sample zone
	// does change the hash
	data = bytes.Repeat([]byte{'A'}, size)
	data[defaultSampleSize-1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h3, _ := HashFilename(sampleFile)
	is.NotEqual(h1, h3)

	data = bytes.Repeat([]byte{'A'}, size)
	data[defaultSampleSize+1] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h4, _ := HashFilename(sampleFile)
	is.NotEqual(h1, h4)
	is.NotEqual(h3, h4)

	data = bytes.Repeat([]byte{'A'}, size)
	data[2*defaultSampleSize] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h5, _ := HashFilename(sampleFile)
	is.NotEqual(h1, h5)
	is.NotEqual(h3, h5)
	is.NotEqual(h4, h5)

	data = bytes.Repeat([]byte{'A'}, size)
	data[2*defaultSampleSize+2] = 'B'
	ioutil.WriteFile(sampleFile, data, 0666)
	h6, _ := HashFilename(sampleFile)
	is.NotEqual(h1, h6)
	is.NotEqual(h3, h6)
	is.NotEqual(h4, h6)
	is.NotEqual(h5, h6)

	// test that changing the size changes the hash
	data = bytes.Repeat([]byte{'A'}, size+1)
	ioutil.WriteFile(sampleFile, data, 0666)
	h7, _ := HashFilename(sampleFile)
	is.NotEqual(h1, h7)
	is.NotEqual(h3, h7)
	is.NotEqual(h4, h7)
	is.NotEqual(h5, h7)
	is.NotEqual(h6, h7)

	os.Remove(sampleFile)
}

func WriteSample(name string, size int) string {
	fullFilename := filepath.Join(tempDir, name)
	data := make([]byte, size)

	for i := 0; i < size; i++ {
		data[i] = 'A'
	}

	err := ioutil.WriteFile(fullFilename, data, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return fullFilename
}
