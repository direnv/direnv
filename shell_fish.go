package main

import (
	"fmt"
	"strings"
)

type fish int

var FISH fish

const FISH_HOOK = `
function __direnv_export_eval --on-event fish_prompt;
	eval (direnv export fish);
end
`

func (f fish) Hook() (string, error) {
	return FISH_HOOK, nil
}

func (f fish) Export(e ShellExport) (out string) {
	for key, value := range e {
		if value == nil {
			out += f.unset(key)
		} else {
			out += f.export(key, *value)
		}
	}
	return out
}

func (f fish) export(key, value string) string {
	if key == "PATH" {
		command := "set -x -g PATH"
		for _, path := range strings.Split(value, ":") {
			command += " " + f.escape(path)
		}
		return command + ";"
	}
	return "set -x -g " + f.escape(key) + " " + f.escape(value) + ";"
}

func (f fish) unset(key string) string {
	return "set -e -g " + f.escape(key) + ";"
}

func (f fish) escape(str string) string {
	in := []byte(str)
	out := "'"
	i := 0
	l := len(in)

	hex := func(char byte) {
		out += fmt.Sprintf("'\\x%02x'", char)
	}

	backslash := func(char byte) {
		out += string([]byte{BACKSLASH, char})
	}

	escaped := func(str string) {
		out += "'" + str + "'"
	}

	literal := func(char byte) {
		out += string([]byte{char})
	}

	for i < l {
		char := in[i]
		switch {
		case char == TAB:
			escaped(`\t`)
		case char == LF:
			escaped(`\n`)
		case char == CR:
			escaped(`\r`)
		case char <= US:
			hex(char)
		case char == SINGLE_QUOTE:
			backslash(char)
		case char == BACKSLASH:
			backslash(char)
		case char <= TILDA:
			literal(char)
		case char == DEL:
			hex(char)
		default:
			hex(char)
		}
		i += 1
	}

	out += "'"

	return out
}
