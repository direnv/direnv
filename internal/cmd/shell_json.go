package cmd

import (
	"encoding/json"
	"errors"
)

// jsonShell is not a real shell
type jsonShell struct{}

// JSON is not really a shell but it fits. Useful to add support to editor and
// other external tools that understand JSON as a format.
var JSON Shell = jsonShell{}

func (sh jsonShell) Hook() (string, error) {
	return "", errors.New("this feature is not supported")
}

func (sh jsonShell) Export(e ShellExport) (string, error) {
	out, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (sh jsonShell) Dump(env Env) (string, error) {
	out, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
