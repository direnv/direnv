package main

import (
	"bytes"
	"encoding/json"
)

type elvish struct{}

var ELVISH = elvish{}

func (elvish) Hook() (string, error) {
	return `## hook for direnv
@edit:before-readline = $@edit:before-readline {
	try {
		m = ("{{.SelfPath}}" export elvish | from-json)
		keys $m | each [k]{
			if (==s $k 'null') {
				unset-env $k
			} else {
				set-env $k $m[$k]
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

var (
	_ Shell = (*elvish)(nil)
)
