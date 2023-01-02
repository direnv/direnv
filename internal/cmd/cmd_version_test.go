package cmd

import (
	"golang.org/x/mod/semver"
	"os"
	"strings"
	"testing"
)

func TestVersionDotTxt(t *testing.T) {
	bs, err := os.ReadFile("../../version.txt")
	if err != nil {
		t.Fatalf("failed to read ../../version.txt: %v", err)
	}
	version = strings.TrimSpace(string(bs))

	if !semver.IsValid(ensureVPrefixed(string(version))) {
		t.Fatalf(`version.txt does not contain a valid semantic version: %q`, version)
	}
}
