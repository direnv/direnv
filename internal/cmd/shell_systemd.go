package cmd

import (
	"errors"
	"strings"
)

// systemdShell is not a real shell
type systemdShell struct{}

// Systemd is not really a shell but is useful to add support
// to systemd EnvironmentFile(https://0pointer.de/public/systemd-man/systemd.exec.html#EnvironmentFile=)
var Systemd Shell = systemdShell{}

func (sh systemdShell) Hook() (string, error) {
	return "", errors.New("this feature is not supported")
}

func (sh systemdShell) Export(e ShellExport) (out string) {
	for key, value := range e {
		if value != nil {
			out += sh.export(key, *value)
		}
	}
	return out
}

func (sh systemdShell) Dump(env Env) (out string) {
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out
}

func logQuoteIssue(key string) {
	// Value contains both type of quotes, meaning we cannot quote around the value
	// as it would make systemd remove that type of quote from the value
	// Log an error as we cannot keep the integrity of the value
	logMsg(`Direnv isn't able to ensure the integrity of the value for key: %v
This is caused by the use of single quotes and double quotes in the value.`, key)
}

func sanitizeValue(key, value string) string {
	var containSpecialChar bool
	specialCharacterList := []string{"\n", "\\"}
	for _, specialChar := range specialCharacterList {
		if strings.ContainsAny(value, specialChar) {
			containSpecialChar = true
		}
	}

	sanitizedValue := value

	if containSpecialChar {
		// Since the value contains special characters it needs to be quoted
		var startsWithDoubleQuotes, startsWithSingleQuotes bool

		_, startsWithDoubleQuotes = strings.CutPrefix(value, "\"")
		_, startsWithSingleQuotes = strings.CutPrefix(value, "'")
		if startsWithDoubleQuotes {
			if strings.ContainsAny(value, "'") {
				logQuoteIssue(key)
			} else {
				// encapsulate with single quotes to preserve all double quotes
				sanitizedValue = "'" + value + "'"
			}

		}
		if startsWithSingleQuotes {
			if strings.ContainsAny(value, "\"") {
				logQuoteIssue(key)
			} else {
				// encapsulate with double quotes to preserve all single quotes
				sanitizedValue = "\"" + value + "\""
			}
		}

	}
	// if the value doesn't contains special characters then we don't touch it
	return sanitizedValue
}

func (sh systemdShell) export(key, value string) string {
	return key + "=" + "\"" + sanitizeValue(key, value) + "\"\n"
}
