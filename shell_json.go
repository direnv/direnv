package main

import (
	"encoding/json"
	"errors"
)

// jsonShell is not a real shell
type jsonShell struct{}

var JSON Shell = jsonShell{}

func (sh jsonShell) Hook() (string, error) {
	return "", errors.New("this feature is not supported")
}

func (sh jsonShell) Export(e ShellExport, q ShellQuotes) string {
	out, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		// Should never happen
		panic(err)
	}
	return string(out)
}

func (sh jsonShell) Dump(env Env) string {
	out, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		// Should never happen
		panic(err)
	}
	return string(out)
}
