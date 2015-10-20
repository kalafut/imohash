package imohash

import (
	"bytes"
	"encoding/binary"
	"hash"
	"io"
	"os"

	"github.com/spaolacci/murmur3"
)

const Size = 16

// Files smaller than 128kb will be hashed in their entirety.
const SampleThreshhold = 128 * 1024
const SampleSize = 16 * 1024

var emptyArray = [Size]byte{}

// Make sure interfaces are correctly implemented.
var (
	_ hash.Hash = new(ImoHash)
)

type ImoHash struct {
	hasher     murmur3.Hash128
	sampleSize int
	bytesAdded int
}

func New(sampleSize ...int) ImoHash {
	h := ImoHash{
		hasher:     murmur3.New128(),
		sampleSize: SampleSize,
	}

	if len(sampleSize) > 0 {
		h.sampleSize = sampleSize[0]
	}

	return h
}

func SumFile(filename string) ([Size]byte, error) {
	h := New(SampleSize)
	return h.SumFile(filename)
}

func (imo *ImoHash) SumFile(filename string) ([Size]byte, error) {
	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return emptyArray, err
	}

	fi, err := f.Stat()
	if err != nil {
		return emptyArray, err
	}
	sr := io.NewSectionReader(f, 0, fi.Size())
	return imo.hashCore(sr), nil
}

func Sum(data []byte) [Size]byte {
	hasher := New(SampleSize)
	sr := io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data)))

	return hasher.hashCore(sr)
}

// hash.Hash methods
func (imo *ImoHash) BlockSize() int { return 1 }

func (imo *ImoHash) Reset() {
	imo.bytesAdded = 0
	imo.hasher.Reset()
}

func (imo *ImoHash) Size() int { return Size }

func (imo *ImoHash) Sum(data []byte) []byte {
	hash := imo.hasher.Sum(nil)
	binary.PutUvarint(hash, uint64(imo.bytesAdded))
	return append(data, hash...)
}

func (imo *ImoHash) Write(data []byte) (n int, err error) {
	imo.hasher.Write(data)
	imo.bytesAdded += len(data)
	return len(data), nil
}

func (imo *ImoHash) hashCore(f *io.SectionReader) [Size]byte {
	var result [Size]byte

	imo.hasher.Reset()

	if f.Size() < SampleThreshhold {
		buffer := make([]byte, f.Size())
		f.Read(buffer)
		imo.hasher.Write(buffer)
	} else {
		buffer := make([]byte, imo.sampleSize)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(f.Size()/2, 0)
		f.Read(buffer)
		imo.hasher.Write(buffer)
		f.Seek(int64(-imo.sampleSize), 2)
		f.Read(buffer)
		imo.hasher.Write(buffer)
	}

	hash := imo.hasher.Sum(nil)

	binary.PutUvarint(hash, uint64(f.Size()))
	copy(result[:], hash)

	return result
}
