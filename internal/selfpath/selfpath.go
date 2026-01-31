// Package selfpath is for computing the utility's own path, required
// for cases where the utility is called as an option or subcommand of
// another program (e.g. `stave --direnv`).
package selfpath

import (
	"regexp"
)

var firstFlagRegexp = regexp.MustCompile(`\s+-`)

// SelfPath processes the original zero argument to handle flag separation and quoting.
func SelfPath(origZeroArg string) string {
	// Find position of first flag (if any)
	firstMatchIndices := firstFlagRegexp.FindStringIndex(origZeroArg)
	if firstMatchIndices == nil {
		return doubleQuote(origZeroArg)
	}

	firstFlagPos := firstMatchIndices[0]

	return doubleQuote(origZeroArg[:firstFlagPos]) + origZeroArg[firstFlagPos:]
}

func doubleQuote(s string) string {
	return `"` + s + `"`
}
