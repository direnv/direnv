package main

import (
	"encoding/json"
	"errors"
)

type jsonShell int

var JSON jsonShell

func (s jsonShell) Hook() (string, error) {
	return "", errors.New("this feature is not supported")
}

func (s jsonShell) Export(e ShellExport) string {
	out, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		// Should never happen
		panic(err)
	}
	return string(out)
}
