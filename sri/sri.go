// Package sri implements helper functions to calculate SubResource Integrity
// hashes.
// https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
package sri

import (
	"crypto/sha256"
	"crypto/sha512"
	b64 "encoding/base64"
	"fmt"
	"hash"
	"io"
	"strings"
)

// Algo is a supported hashing algorithm
type Algo string

const (
	// SHA256 algo
	SHA256 = Algo("sha256")
	// SHA384 algo
	SHA384 = Algo("sha384")
	// SHA512 algo
	SHA512 = Algo("sha512")
)

// GetAlgo extracts the algo part of the hash
func GetAlgo(sriHash string) (Algo, error) {
	elems := strings.SplitN(sriHash, "-", 2)
	if len(elems) != 2 {
		return Algo(""), fmt.Errorf("not a SRI hash %v", sriHash)
	}
	switch elems[0] {
	case string(SHA256):
		return SHA256, nil
	case string(SHA384):
		return SHA384, nil
	case string(SHA512):
		return SHA512, nil
	default:
		return Algo(""), fmt.Errorf("unsupported SRI also %s", elems[0])
	}
}

// Writer is like a hash.Hash with a Sum function
type Writer struct {
	w    io.Writer
	algo Algo
	h    hash.Hash
}

// NewWriter returns a SRI writer that forwards the write while calculating
// the SRI hash.
func NewWriter(w io.Writer, algo Algo) Writer {
	var h hash.Hash
	switch algo {
	case SHA256:
		h = sha256.New()
	case SHA384:
		h = sha512.New384()
	case SHA512:
		h = sha512.New()
	default:
		panic("unsupported SRI algo")
	}

	return Writer{w, algo, h}
}

func (w Writer) Write(b []byte) (int, error) {
	// First write to the underlying storage
	n, err := w.w.Write(b)
	if err == nil {
		// This should always succeed
		_, _ = w.h.Write(b)
	}
	return n, err
}

// Sum returns the calculated SRI hash
func (w Writer) Sum() string {
	hashResult := w.h.Sum(nil)
	b64Result := b64.StdEncoding.EncodeToString(hashResult)
	return string(w.algo) + "-" + b64Result
}
