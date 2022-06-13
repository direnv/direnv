package sri

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
)

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
func (w Writer) Sum() *Hash {
	sum := w.h.Sum(nil)
	return &Hash{string(w.algo), sum}
}
