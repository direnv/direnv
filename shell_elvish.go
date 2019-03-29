package main

import (
	"bytes"
	"encoding/json"
)

type elvish struct{}

var ELVISH Shell = elvish{}

func (elvish) Hook() (string, error) {
	return `## hook for direnv
@edit:before-readline = $@edit:before-readline {
	try {
		m = [("{{.SelfPath}}" export elvish | from-json)]
		if (> (count $m) 0) {
			m = (explode $m)
			keys $m | each [k]{
				if (==s $k 'null') {
					unset-env $k
				} else {
					set-env $k $m[$k]
				}
			}
		}
	} except e {
		echo $e
	}
}
`, nil
}

func (sh elvish) Export(e ShellExport) string {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(e)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (sh elvish) Dump(env Env) (out string) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(env)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

var (
	_ Shell = (*elvish)(nil)
)
