package cmd

import (
	"bytes"
	"encoding/json"
)

type murex struct{}

// Murex is the shell implementation for Murex shell.
var Murex Shell = murex{}

const murexHook = `event: onPrompt direnv_hook=before {
	"{{.SelfPath}}" export murex -> set exports
	if { $exports != "" } {
		$exports -> :json: formap key value {
			if { is-null value } then {
				!export "$key"
			} else {
				$value -> export "$key"
			}
		}
	}
}`

func (sh murex) Hook() (string, error) {
	return murexHook, nil
}

func (sh murex) Dump(env Env) (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(env)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (sh murex) Export(e ShellExport) (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(e)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

var (
	_ Shell = (*murex)(nil)
)
