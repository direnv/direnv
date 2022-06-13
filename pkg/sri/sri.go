// Package sri implements helper functions to calculate SubResource Integrity
// hashes.
// https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
package sri

import (
	b64 "encoding/base64"
	"encoding/hex"
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

// Base64 encoding to use
var b64Enc = b64.StdEncoding

// Hash represents a SRI-hash
type Hash struct {
	algo string
	sum  []byte
}

// String returns a SRI-encoded string
func (h *Hash) String() string {
	return h.algo + "-" + b64Enc.EncodeToString(h.sum)
}

// Hex return a hex-encoded representation of the sum
func (h *Hash) Hex() string {
	return hex.EncodeToString(h.sum)
}
