package main

import (
	"testing"
)

func TestShellEscape(t *testing.T) {
	assertEqual(t, "''", ShellEscape(""))
	assertEqual(t, `'\n'`, ShellEscape("\n"))
	assertEqual(t, `\ `, ShellEscape(" "))
}

func assertEqual(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("Expected \"%v\" to equal \"%v\"", b, a)
	}
}
