package imohash

import (
	"crypto/md5"
	"fmt"
	"testing"

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
		{16384, 131072, 1, "01659e2ec0f3c75bf39e43a41adb5d4f"},
		{16384, 131072, 127, "7f47671cc79d4374404b807249f3166e"},
		{16384, 131072, 128, "800183e5dbea2e5199ef7c8ea963a463"},
		{16384, 131072, 4095, "ff1f770d90d3773949d89880efa17e60"},
		{16384, 131072, 4096, "802048c26d66de432dbfc71afca6705d"},
		{16384, 131072, 131072, "8080085a3d3af2cb4b3a957811cdf370"},
		{16384, 131073, 131072, "808008282d3f3b53e1fd132cc51fcc1d"},
		{16384, 131072, 500000, "a0c21e44a0ba3bddee802a9d1c5332ca"},
		{50, 131072, 300000, "e0a712edd8815c606344aed13c44adcf"},
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
	hasher := md5.New()

	for len(r) < n {
		hasher.Write([]byte{'A'})
		r = hasher.Sum(r)
	}

	return r[0:n]
}
