package sri

import (
	"fmt"
	"strings"
)

// Parse a SRI hash
func Parse(sriHash string) (*Hash, error) {
	elems := strings.SplitN(sriHash, "-", 2)
	if len(elems) != 2 {
		return nil, fmt.Errorf("sri: not a hash %v", sriHash)
	}

	// Get the algo
	var algo Algo
	switch elems[0] {
	case string(SHA256):
		algo = SHA256
	case string(SHA384):
		algo = SHA384
	case string(SHA512):
		algo = SHA512
	default:
		return nil, fmt.Errorf("sri: unsupported algo %s", elems[0])
	}

	// Get the hash
	dbuf := make([]byte, b64Enc.DecodedLen(len(elems[1])))
	n, err := b64Enc.Decode(dbuf, []byte(elems[1]))
	if err != nil {
		return nil, err
	}
	sum := dbuf[:n]

	return &Hash{string(algo), sum}, nil
}
