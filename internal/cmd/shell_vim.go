package cmd

import (
	"errors"
	"strings"
)

type vim struct{}

// Vim adds support for vim. Not really a shell but it's handly.
var Vim Shell = vim{}

func (sh vim) Hook() (string, error) {
	return "", errors.New("this feature is not supported. Install the direnv.vim plugin instead")
}

func (sh vim) Export(e ShellExport) (out string) {
	for key, value := range e {
		if value == nil {
			out += sh.unset(key)
		} else {
			out += sh.export(key, *value)
		}
	}
	return out
}

func (sh vim) Dump(env Env) (out string) {
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out
}

func (sh vim) export(key, value string) string {
	return "call setenv(" + sh.escapeKey(key) + "," + sh.escapeValue(value) + ")\n"
}

func (sh vim) unset(key string) string {
	return "call setenv(" + sh.escapeKey(key) + ",v:null)\n"
}

// TODO: support keys with special chars or fail
func (sh vim) escapeKey(str string) string {
	return sh.escapeValue(str)
}

// TODO: Make sure this escaping is valid
func (sh vim) escapeValue(str string) string {
	return "'" + strings.Replace(str, "'", "''", -1) + "'"
}
