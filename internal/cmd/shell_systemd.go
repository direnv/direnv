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

func cutEncapsulated(valueToTest, encapsulatingValue string) (cutValue string, wasEncapsulated bool) {
	withoutPrefix, startsWithEncapsulatingValue := strings.CutPrefix(valueToTest, encapsulatingValue)
	if startsWithEncapsulatingValue {
		withoutPrefixAndSuffix, endsWithEncapsulatingValue := strings.CutSuffix(withoutPrefix, encapsulatingValue)
		if endsWithEncapsulatingValue {
			return withoutPrefixAndSuffix, true
		}
	}
	return valueToTest, false
}

func sanitizeValue(value string) string {
	containSpecialChar := false
	specialCharacterList := []string{"\n", "\\", `"`, `'`}
	for _, specialChar := range specialCharacterList {
		if strings.ContainsAny(value, specialChar) {
			containSpecialChar = true
		}
	}

	sanitizedValue := value

	if containSpecialChar {
		// Since the value contains special characters it needs to be quoted

		valueWithoutEncapsulation, encapsulatedBySingleQuotes := cutEncapsulated(value, `'`)

		if encapsulatedBySingleQuotes {
			sanitizedValue = `'` + strings.ReplaceAll(valueWithoutEncapsulation, `'`, `\'`) + `'`
		} else {
			valueWithoutEncapsulation, encapsulatedByDoubleQuotes := cutEncapsulated(value, `"`)
			logDebug("encapsulated by double quotes : %v", encapsulatedByDoubleQuotes)
			sanitizedValue = `"` + strings.ReplaceAll(valueWithoutEncapsulation, `"`, `\"`) + `"`
		}
	}
	// if the value doesn't contains special characters then we don't touch it
	return sanitizedValue
}

func (sh systemdShell) export(key, value string) string {
	return key + "=" + sanitizeValue(value) + "\n"
}
