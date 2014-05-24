package main

import (
	"fmt"
	"path/filepath"
)

/*
 * Shells
 */
type Shell interface {
	Hook() string
	Export(key, value string) string
	Unset(key string) string
}

func DetectShell(target string) Shell {
	target = filepath.Base(target)
	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	switch target {
	case "bash":
		return BASH
	case "zsh":
		return ZSH
	case "fish":
		return FISH
	case "vim":
		return VIM
	}

	return nil
}

/*
 * Escaping
 */

const (
	ACK           = 6
	TAB           = 9
	LF            = 10
	CR            = 13
	US            = 31
	AMPERSTAND    = 38
	SINGLE_QUOTE  = 39
	PLUS          = 43
	NINE          = 57
	QUESTION      = 63
	LOWERCASE_Z   = 90
	OPEN_BRACKET  = 91
	BACKSLASH     = 92
	UNDERSCORE    = 95
	CLOSE_BRACKET = 93
	BACKTICK      = 96
	TILDA         = 126
	DEL           = 127
)

// https://github.com/solidsnack/shell-escape/blob/master/Text/ShellEscape/Bash.hs
/*
A Bash escaped string. The strings are wrapped in @$\'...\'@ if any
bytes within them must be escaped; otherwise, they are left as is.
Newlines and other control characters are represented as ANSI escape
sequences. High bytes are represented as hex codes. Thus Bash escaped
strings will always fit on one line and never contain non-ASCII bytes.
*/
func ShellEscape(str string) string {
	if str == "" {
		return "''"
	}
	in := []byte(str)
	out := ""
	i := 0
	l := len(in)
	escape := false

	hex := func(char byte) {
		escape = true
		out += fmt.Sprintf("\\x%02x", char)
	}

	backslash := func(char byte) {
		escape = true
		out += string([]byte{BACKSLASH, char})
	}

	escaped := func(str string) {
		escape = true
		out += str
	}

	quoted := func(char byte) {
		escape = true
		out += string([]byte{char})
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
		case char <= LOWERCASE_Z:
			literal(char)
		case char == OPEN_BRACKET:
			quoted(char)
		case char == BACKSLASH:
			backslash(char)
		case char <= CLOSE_BRACKET:
			quoted(char)
		case char == UNDERSCORE:
			literal(char)
		case char <= BACKTICK:
			quoted(char)
		case char <= LOWERCASE_Z:
			literal(char)
		case char <= TILDA:
			quoted(char)
		case char == DEL:
			hex(char)
		default:
			hex(char)
		}
		i += 1
	}

	if escape {
		out = "$'" + out + "'"
	}

	return out
}
