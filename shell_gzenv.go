package main

import (
	"errors"

	"github.com/direnv/direnv/gzenv"
)

// gzenvShell is not a real shell. used for internal purposes.
type gzenvShell int

var GZENV Shell = gzenvShell(0)

func (s gzenvShell) Hook() (string, error) {
	return "", errors.New("the gzenv shell doesn't support hooking")
}

func (s gzenvShell) Export(e ShellExport, q ShellQuotes) string {
	return gzenv.Marshal(e)
}

func (s gzenvShell) Dump(env Env) string {
	return gzenv.Marshal(env)
}
