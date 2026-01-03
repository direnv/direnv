package selfpath

import (
	"regexp"
)

var firstFlagRegexp = regexp.MustCompile(`\s+-`)

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
