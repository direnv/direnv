package shell

import (
	"regexp"
	"strings"
)

// This function and comments have been copied over from Ruby's
// stdlib shellwords.rb library.
func Escape(str string) string {
	if str == "" {
		return "''"
	}

	// Treat multibyte characters as is.  It is caller's responsibility
	// to encode the string in the right encoding for the shell
	// environment.
	r := regexp.MustCompile("([^A-Za-z0-9_\\-.,:/@\n])")
	replace := func(match string) string { return "\\\\\\" + match }
	str = r.ReplaceAllStringFunc(str, replace)

	// A LF cannot be escaped with a backslash because a backslash + LF
	// combo is regarded as line continuation and simply ignored.
	str = strings.Replace(str, "\n", "'\n'", -1)
	return str
}

