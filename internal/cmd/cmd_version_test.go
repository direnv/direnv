package cmd

import (
	"golang.org/x/mod/semver"
	"io/ioutil"
	"testing"
)

func TestVersionDotTxt(t *testing.T) {
	version, _ := ioutil.ReadFile("../../version.txt")

	if !semver.IsValid(ensureVPrefixed(string(version))) {
		t.Fatalf(`version.txt does not contain a valid semantic version: %q`, version)
	}
}
