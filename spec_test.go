package imohash

import (
	"fmt"
	"testing"

	"github.com/spaolacci/murmur3"

	"gopkg.in/tylerb/is.v1"
)

func TestSpec(t *testing.T) {
	var hashStr string
	is := is.New(t)

	tests := []struct {
		s    int
		t    int
		n    int
		hash string
	}{
		{16384, 131072, 0, "00000000000000000000000000000000"},
		{16384, 131072, 1, "016ac6dd306a3e594e711127c5b5a8e4"},
		{16384, 131072, 127, "7f0e0eaa0a415e2878303214a087fbab"},
		{16384, 131072, 128, "80017902d85ab752459758292a217ed8"},
		{16384, 131072, 4095, "ff1fc7e3b9890469ee676df99f6ac7b7"},
		{16384, 131072, 4096, "80203d4ea589349dcdf63511a8d49d4a"},
		{16384, 131072, 131072, "808008ed886018d5cd37a9b35bbae286"},
		{16384, 131073, 131072, "808008ac6ab33ddcaadea68ba5ad4f05"},
		{16384, 131072, 500000, "a0c21ea46738cf366d496962000a45f7"},
	}

	for _, test := range tests {
		i := NewCustom(test.s, test.t)
		hashStr = fmt.Sprintf("%x", i.Sum(M(test.n)))
		is.Equal(hashStr, test.hash)
	}
}

// M generates n bytes of pseudo-random data according to the
// method described in the imohash algorithm description.
func M(n int) []byte {
	r := make([]byte, 0, n)
	hasher := murmur3.New128()

	for len(r) < n {
		hasher.Write([]byte{'A'})
		r = hasher.Sum(r)
	}

	return r[0:n]
}
