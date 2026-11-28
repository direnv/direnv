package cmd

import (
	"strings"
)

type tcsh struct{}

// Tcsh adds support for the tickle shell
var Tcsh Shell = tcsh{}

func (sh tcsh) Hook() (string, error) {
	return "alias precmd 'eval `{{.SelfPath}} export tcsh`'", nil
}

func (sh tcsh) Export(e ShellExport) (string, error) {
	var out strings.Builder
	for key, value := range e {
		if value == nil {
			out.WriteString(sh.unset(key))
		} else {
			out.WriteString(sh.export(key, *value))
		}
	}
	return out.String(), nil
}

func (sh tcsh) Dump(env Env) (string, error) {
	var out strings.Builder
	for key, value := range env {
		out.WriteString(sh.export(key, value))
	}
	return out.String(), nil
}

func (sh tcsh) export(key, value string) string {
	if key == "PATH" {
		var command strings.Builder
		command.WriteString("set path = (")
		for path := range strings.SplitSeq(value, ":") {
			command.WriteString(" " + sh.escape(path))
		}
		return command.String() + " );"
	}
	return "setenv " + sh.escape(key) + " " + sh.escape(value) + " ;"
}

func (sh tcsh) unset(key string) string {
	return "unsetenv " + sh.escape(key) + " ;"
}

func (sh tcsh) escape(str string) string {
	if str == "" {
		return "''"
	}
	in := []byte(str)
	var out strings.Builder
	out.Grow(len(in))
	i := 0
	l := len(in)

	backslash := func(char byte) {
		out.WriteByte(BACKSLASH)
		out.WriteByte(char)
	}

	// tcsh interprets \t, \n and \r inside double quotes.
	escaped := func(str string) {
		out.WriteByte(DOUBLE_QUOTE)
		out.WriteString(str)
		out.WriteByte(DOUBLE_QUOTE)
	}

	quoted := func(char byte) {
		// Raw bytes inside double quotes pass through tcsh untouched; it is
		// the only reliable way to emit control and non-ASCII characters.
		out.WriteByte(DOUBLE_QUOTE)
		out.WriteByte(char)
		out.WriteByte(DOUBLE_QUOTE)
	}

	literal := func(char byte) {
		out.WriteByte(char)
	}

	for i < l {
		char := in[i]
		switch {
		case char == ACK:
			quoted(char)
		case char == TAB:
			escaped(`\t`)
		case char == LF:
			escaped(`\n`)
		case char == CR:
			escaped(`\r`)
		case char == SPACE:
			backslash(char)
		case char == DOUBLE_QUOTE:
			backslash(char)
		case char <= US:
			quoted(char)
		case char <= AMPERSTAND:
			quoted(char)
		case char == SINGLE_QUOTE:
			backslash(char)
		case char <= PLUS:
			quoted(char)
		case char <= NINE:
			literal(char)
		case char <= QUESTION:
			quoted(char)
		case char <= UPPERCASE_Z:
			literal(char)
		case char == OPEN_BRACKET:
			quoted(char)
		case char == BACKSLASH:
			backslash(char)
		case char == UNDERSCORE:
			literal(char)
		case char <= LOWERCASE_Z:
			literal(char)
		case char <= CLOSE_BRACKET:
			quoted(char)
		case char <= BACKTICK:
			quoted(char)
		case char <= TILDE:
			quoted(char)
		case char == DEL:
			quoted(char)
		default:
			quoted(char)
		}
		i++
	}

	return out.String()
}
