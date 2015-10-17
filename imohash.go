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

func HashFilename(filename string) (uint64, error) {
	return defaultHasher.HashFilename(filename)
}

func HashFile(file *os.File) (uint64, error) {
	return defaultHasher.HashFile(file)
}

func (imo ImoHash) HashFilename(file string) (uint64, error) {
	f, err := os.Open(file)
	defer f.Close()

	if err != nil {
		return 0, err
	}

	return imo.HashFile(f)
}

func (imo ImoHash) HashFile(f *os.File) (uint64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	imo.hasher.Reset()

	if fi.Size() <= imo.sampleSize*3 {
		buffer := make([]byte, fi.Size())

		f.Read(buffer)
		imo.hasher.Write(buffer)
	} else {
		buffer := make([]byte, imo.sampleSize)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(fi.Size()/2-imo.sampleSize/2, 0)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(-imo.sampleSize, 2)
		f.Read(buffer)
		imo.hasher.Write(buffer)
	}

	size := foldInt(fi.Size())

	return (uint64(size) << 32) | uint64(imo.hasher.Sum32()), nil
}

func foldInt(v int64) uint32 {
	var r, i uint32

	for i = 0; i < 4; i++ {
		r |= uint32((byte(v>>(8*i)) ^ byte(v>>(56-(8*i))))) << (8 * i)
	}

	return r
}
