package main

import (
	"fmt"
	"strings"
)

type fish int

var FISH fish

func (f fish) Hook() string {
	return `
function __direnv_export_eval --on-event fish_prompt;
	eval (direnv export fish);
end
`
}

func (f fish) Escape(str string) string {
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

func (f fish) Export(key, value string) string {
	if key == "PATH" {
		command := "set -x -g PATH"
		for _, path := range strings.Split(value, ":") {
			command += " " + f.Escape(path)
		}
		return command + ";"
	}
	return "set -x -g " + f.Escape(key) + " " + f.Escape(value) + ";"
}

func (f fish) Unset(key string) string {
	return "set -e -g " + f.Escape(key) + ";"
}
