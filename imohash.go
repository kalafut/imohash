package imohash

import (
	"hash"
	"os"

	"github.com/spaolacci/murmur3"
)

const defaultSampleSize = 4096
const FULL_HASH_LIMIT = 3 * defaultSampleSize

var defaultHasher = New(defaultSampleSize)

type ImoHash struct {
	hasher     hash.Hash32
	sampleSize int64
}

func New(sampleSize int64) ImoHash {
	if sampleSize < 1 {
		sampleSize = defaultSampleSize
	}

	h := ImoHash{
		hasher:     murmur3.New32(),
		sampleSize: sampleSize,
	}

	return h
}

func HashFilename(filename string) (uint32, error) {
	return defaultHasher.HashFilename(filename)
}

func HashFile(file *os.File) (uint32, error) {
	return defaultHasher.HashFile(file)
}

func (imo ImoHash) HashFilename(file string) (uint32, error) {
	f, err := os.Open(file)
	defer f.Close()

	if err != nil {
		return 0, err
	}

	return imo.HashFile(f)
}

func (imo ImoHash) HashFile(f *os.File) (uint32, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	imo.hasher.Reset()

	imo.hasher.Write(intToSlice(fi.Size()))
	if fi.Size() <= imo.sampleSize*3 {
		buffer := make([]byte, fi.Size())

		f.Read(buffer)
		imo.hasher.Write(buffer)
	} else {
		buffer := make([]byte, imo.sampleSize)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(fi.Size()/2, 0)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(-imo.sampleSize, 2)
		f.Read(buffer)
		imo.hasher.Write(buffer)
	}

	return imo.hasher.Sum32(), nil
}

func intToSlice(i int64) []byte {
	return []byte{
		byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32),
		byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i),
	}
}
