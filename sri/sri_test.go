package sri

import (
	"strings"
	"testing"
)

func TestSRIHasher(t *testing.T) {
	var b strings.Builder

	s := "testdata"

	// Generated with:
	// `echo -n "testdata" | openssl dgst -sha256 -binary - | // openssl base64 -A`
	expectedHash := "sha256-gQ/y+yQqXe5CIPLLDmpRmJH7Z/L4KKbKtO+IlGM7H1A="

	w := NewWriter(&b, SHA256)

	// Check the writer
	n, err := w.Write([]byte(s))
	if err != nil {
		t.Fatalf("write error: %s", err)
	}
	if n != len(s) {
		t.Fatalf("expected len %d but got %d", len(s), n)
	}
	if b.String() != s {
		t.Fatal("data has not been forwarded")
	}

	// Check that the hash has been calculated properly
	x := w.Sum()
	if x != expectedHash {
		t.Fatal("hash mismatch")
	}
}
