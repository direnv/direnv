package cmd

import (
	"errors"
)

// environShell is not a real shell
type environShell struct{}

// Environ is a format exporter and not a shell. It implements a similar
// format as /proc/<pid>/environ in Linux.
var Environ Shell = environShell{}

func (sh environShell) Hook() (string, error) {
	return "", errors.New("this feature is not supported")
}

// Exports emits a string composed of <key>(=<value>), separated by \0.
//
// It's the same as the Dump format with one additional case: if there is no =
// in the line, the environment variable should be unset.
func (sh environShell) Export(e ShellExport) (out string) {
	for key, value := range e {
		if value == nil {
			out += sh.unset(key)
		} else {
			out += sh.export(key, *value)
		}
	}
	return out
}

// Dump emits a string composed of <key>=<value, separated by \0. That
// format is the same as you would find in /prod/<pid>/environ and works with
// values that include special characters like \n.
func (sh environShell) Dump(env Env) (out string) {
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out
}

// <key>=<value> , terminated by \0
func (sh environShell) export(key, value string) string {
	return key + "=" + value + string([]byte{0})
}

// <key> , terminated by \0
func (sh environShell) unset(key string) string {
	return key + string([]byte{0})
}
