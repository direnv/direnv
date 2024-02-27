package cmd

import (
	"bytes"
	"encoding/json"
)

type murex struct{}

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

func (sh murex) Dump(env Env) (out string) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(env)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (sh murex) Export(e ShellExport) string {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(e)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

var (
	_ Shell = (*murex)(nil)
)
