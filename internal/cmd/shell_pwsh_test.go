package cmd

import (
	"strings"
	"testing"
)

// If direnv is installed under a path containing spaces (a common case on
// Windows, e.g. under "C:\Program Files\..."), the hook must invoke it via
// the call operator with a quoted path, or PowerShell splits the path on
// the space and fails to find the command, silently breaking the hook.
func TestPwshHookQuotesSelfPath(t *testing.T) {
	hook, err := Pwsh.Hook()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(hook, `(& "{{.SelfPath}}" export pwsh)`) {
		t.Fatalf("expected the hook to invoke a quoted SelfPath via the call operator, got: %s", hook)
	}
}
