package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDotEnvRejectsUnknownShell(t *testing.T) {
	t.Setenv("PWD", t.TempDir())
	target := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(target, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, shellName := range []string{"unknown", ""} {
		t.Run(shellName, func(t *testing.T) {
			err := cmdDotEnvAction(Env{}, []string{"dotenv", shellName, target})
			want := "unknown target shell '" + shellName + "'"
			if err == nil || err.Error() != want {
				t.Fatalf("expected %q, got %v", want, err)
			}
		})
	}
}
