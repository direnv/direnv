package cmd

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/mod/semver"
)

func TestVersionDotTxt(t *testing.T) {
	bs, err := os.ReadFile("../callable/version.txt")
	if err != nil {
		t.Fatalf("failed to read ../callable/version.txt: %v", err)
	}
	version = strings.TrimSpace(string(bs))

	if !semver.IsValid(ensureVPrefixed(string(version))) {
		t.Fatalf(`version.txt does not contain a valid semantic version: %q`, version)
	}
}
