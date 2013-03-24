package main

import (
	"testing"
)

func TestShellEscape(t *testing.T) {
	assertEqual(t, `''`, ShellEscape(""))
	assertEqual(t, `$'escape\'quote'`, ShellEscape("escape'quote"))
	assertEqual(t, `$'foo\r\n\tbar'`, ShellEscape("foo\r\n\tbar"))
	assertEqual(t, `$'foo bar'`, ShellEscape("foo bar"))
	assertEqual(t, `$'\xc3\xa9'`, ShellEscape("é"))
}

func assertEqual(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("Expected \"%v\" to equal \"%v\"", b, a)
	}
}
