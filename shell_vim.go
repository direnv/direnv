package main

import (
	"os"
	"strings"
)

type vim int

var VIM vim

func (x vim) Hook() string {
	log_error("this feature is not supported. Install the direnv.vim plugin instead.")
	os.Exit(1)
	return ""
}

// TODO: support keys with special chars or fail
func (x vim) EscapeKey(str string) string {
	return str
}

// TODO: Make sure this escaping is valid
func (x vim) EscapeValue(str string) string {
	return "'" + strings.Replace(str, "'", "''", -1) + "'"
}

func (x vim) Export(key, value string) string {
	return "let $" + x.EscapeKey(key) + " = " + x.EscapeValue(value) + "\n"
}

func (x vim) Unset(key string) string {
	return "unlet $" + x.EscapeKey(key) + "\n"
}
