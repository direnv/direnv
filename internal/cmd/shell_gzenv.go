package cmd

import (
	"errors"

	"github.com/direnv/direnv/v2/gzenv"
)

type gzenvShell int

// GzEnv is not a real shell. used for internal purposes.
var GzEnv Shell = gzenvShell(0)

func (s gzenvShell) Hook() (string, error) {
	return "", errors.New("the gzenv shell doesn't support hooking")
}

func (s gzenvShell) Export(e ShellExport) string {
	return gzenv.Marshal(e)
}

func (s gzenvShell) Dump(env Env) string {
	return gzenv.Marshal(env)
}
