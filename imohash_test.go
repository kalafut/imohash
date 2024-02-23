package imohash

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var tempDir string

func TestMain(m *testing.M) {
	flag.Parse()

	// Make a temp area for test files
	tempDir, _ = os.MkdirTemp(os.TempDir(), "imohash_test_data")
	ret := m.Run()
	os.RemoveAll(tempDir)
	os.Exit(ret)
}

func TestCustom(t *testing.T) {
	const sampleFile = "sample"
	var hash [Size]byte
	var err error

	sampleSize := 3
	sampleThreshold := 45
	imo := NewCustom(sampleSize, sampleThreshold)

	// empty file
	os.WriteFile(sampleFile, []byte{}, 0666)
	hash, err = imo.SumFile(sampleFile)
	ok(t, err)
	equal(t, hash, [Size]byte{})

	// small file
	os.WriteFile(sampleFile, []byte("hello"), 0666)
	hash, err = imo.SumFile(sampleFile)
	ok(t, err)

	hashStr := fmt.Sprintf("%x", hash)
	equal(t, hashStr, "05d8a7b341bd9b025b1e906a48ae1d19")

	/* boundary tests using the custom sample size */
	size := sampleThreshold

	// test that changing the gaps between sample zones does not affect the hash
	data := bytes.Repeat([]byte{'A'}, size)
	os.WriteFile(sampleFile, data, 0666)
	h1, _ := imo.SumFile(sampleFile)

	data[sampleSize] = 'B'
	data[size-sampleSize-1] = 'B'
	os.WriteFile(sampleFile, data, 0666)
	h2, _ := imo.SumFile(sampleFile)
	equal(t, h1, h2)

	// test that changing a byte on the edge (but within) a sample zone
	// does change the hash
	data = bytes.Repeat([]byte{'A'}, size)
	data[sampleSize-1] = 'B'
	os.WriteFile(sampleFile, data, 0666)
	h3, _ := imo.SumFile(sampleFile)
	notEqual(t, h1, h3)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size/2] = 'B'
	os.WriteFile(sampleFile, data, 0666)
	h4, _ := imo.SumFile(sampleFile)
	notEqual(t, h1, h4)
	notEqual(t, h3, h4)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size/2+sampleSize-1] = 'B'
	os.WriteFile(sampleFile, data, 0666)
	h5, _ := imo.SumFile(sampleFile)
	notEqual(t, h1, h5)
	notEqual(t, h3, h5)
	notEqual(t, h4, h5)

	data = bytes.Repeat([]byte{'A'}, size)
	data[size-sampleSize] = 'B'
	os.WriteFile(sampleFile, data, 0666)
	h6, _ := imo.SumFile(sampleFile)
	notEqual(t, h1, h6)
	notEqual(t, h3, h6)
	notEqual(t, h4, h6)
	notEqual(t, h5, h6)

	// test that changing the size changes the hash
	data = bytes.Repeat([]byte{'A'}, size+1)
	os.WriteFile(sampleFile, data, 0666)
	h7, _ := imo.SumFile(sampleFile)
	notEqual(t, h1, h7)
	notEqual(t, h3, h7)
	notEqual(t, h4, h7)
	notEqual(t, h5, h7)
	notEqual(t, h6, h7)

	// test sampleSize < 1
	imo = NewCustom(0, size)
	data = bytes.Repeat([]byte{'A'}, size)
	os.WriteFile(sampleFile, data, 0666)
	hash, _ = imo.SumFile(sampleFile)
	hashStr = fmt.Sprintf("%x", hash)
	equal(t, hashStr, "2d9123b54d37e9b8f94ab37a7eca6f40")

	os.Remove(sampleFile)
}

// Test that the top level functions are the same as custom
// functions using the spec defaults.
func TestDefault(t *testing.T) {
	const sampleFile = "sample"
	var h1, h2 [Size]byte
	var testData []byte

	for _, size := range []int{100, 131071, 131072, 50000} {
		imo := NewCustom(16384, 131072)
		testData = M(size)
		equal(t, Sum(testData), imo.Sum(testData))
		os.WriteFile(sampleFile, []byte{}, 0666)
		h1, _ = SumFile(sampleFile)
		h2, _ = imo.SumFile(sampleFile)
		equal(t, h1, h2)
	}
	os.Remove(sampleFile)
}

// Testing helpers from: https://github.com/benbjohnson/testing

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equal fails the test if exp is not equal to act.
func equal(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// equal fails the test if exp is equal to act.
func notEqual(tb testing.TB, exp, act interface{}) {
	if reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texpected mismatch, got matching\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, act)
		tb.FailNow()
	}
}
