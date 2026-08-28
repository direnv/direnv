package cmd

import (
	"strings"
	"testing"
)

// On Windows, `PATH` is inherited from the native process as `Path`, while a
// POSIX shell (Bash) used to evaluate .envrc always normalizes it back to
// `PATH`. This produces a ShellExport with a removal of `Path` and an
// addition of `PATH` for the very same, case-insensitive Windows
// environment variable. Because PowerShell environment variables are
// case-insensitive, the removal must never be allowed to run after the
// addition, or PATH would end up empty.
func TestPwshExportCaseInsensitivePathCollision(t *testing.T) {
	e := make(ShellExport)
	e.Remove("Path")
	e.Add("PATH", "/usr/bin")

	for range 100 {
		out, err := Pwsh.Export(e)
		if err != nil {
			t.Fatal(err)
		}

		if strings.Contains(out, "Remove-Item") {
			t.Fatalf("expected no removal of the case-insensitively colliding PATH variable, got: %s", out)
		}

		if !strings.Contains(out, "${env:PATH}='/usr/bin';") {
			t.Fatalf("expected PATH to be set, got: %s", out)
		}
	}
}

func TestPwshExportUnrelatedRemovalStillHappens(t *testing.T) {
	e := make(ShellExport)
	e.Remove("FOO")
	e.Add("BAR", "baz")

	out, err := Pwsh.Export(e)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "Remove-Item -LiteralPath 'env:/FOO';") {
		t.Fatalf("expected unrelated FOO removal to still happen, got: %s", out)
	}

	if !strings.Contains(out, "${env:BAR}='baz';") {
		t.Fatalf("expected BAR to be set, got: %s", out)
	}
}
