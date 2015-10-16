package imohash

import (
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

func TestIntToSlice(t *testing.T) {
	is := is.New(t)

	is.Equal(intToSlice(0), []byte{0, 0, 0, 0, 0, 0, 0, 0})
	is.Equal(intToSlice(1), []byte{0, 0, 0, 0, 0, 0, 0, 1})
	is.Equal(intToSlice(0x44eeddccbbaa9977), []byte{0x44, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x77})
}

func TestBasic(t *testing.T) {
	var hash uint32
	var err error

	is := is.New(t)

	// empty file
	ioutil.WriteFile("sample", []byte{}, 0666)
	hash, err = HashFilename("sample")
	is.NotErr(err)
	is.Equal(hash, murmur3.Sum32([]byte{0, 0, 0, 0, 0, 0, 0, 0}))

	// small file
	ioutil.WriteFile("sample", []byte("hello"), 0666)
	hash, err = HashFilename("sample")
	is.Equal(hash, murmur3.Sum32([]byte{0, 0, 0, 0, 0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o'}))

	filename := WriteSample("s", 100)

	hash, err = HashFilename(filename)
	is.NotErr(err)
	//is.Equal(hash, 0xce42b64)
}

/*
func TestHash(t *testing.T) {
	SAMPLE_FILE := filepath.Join(tempDir, "sample.txt")

	is := is.New(t)

	h, err := smartHash("not_found.txt")
	is.Err(err)

	WriteSample(SAMPLE_FILE, 100)
	h, err = smartHash(SAMPLE_FILE)
	os.Remove(SAMPLE_FILE)
	is.NotErr(err)

	is.Equal(h, uint64(0xf98540d8f8a71e22))

	WriteSample(SAMPLE_FILE, 10000000)
	h, err = smartHash(SAMPLE_FILE)
	os.Remove(SAMPLE_FILE)
	is.NotErr(err)

	is.Equal(h, uint64(0x93686806171fdb95))

	WriteSample(SAMPLE_FILE, 10000001)
	h, err = smartHash(SAMPLE_FILE)
	os.Remove(SAMPLE_FILE)
	is.NotErr(err)

	is.Equal(h, uint64(0x93686806171fdb95))
}
*/

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
