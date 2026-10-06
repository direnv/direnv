package cmd

import (
	"strings"
	"testing"
	"text/template"
)

func TestPwshExportRemovesBeforeAssigning(t *testing.T) {
	e := ShellExport{}
	e.Remove("Path")
	e.Add("PATH", `C:\bin`)
	e.Remove("SystemRoot")
	e.Add("SYSTEMROOT", `C:\WINDOWS`)

	for i := 0; i < 100; i++ {
		out, err := Pwsh.Export(e)
		if err != nil {
			t.Fatal(err)
		}
		lastRemove := strings.LastIndex(out, "Remove-Item")
		firstAssign := strings.Index(out, "${env:")
		if lastRemove > firstAssign {
			t.Fatalf("assignment emitted before removal: %s", out)
		}
	}
}

func TestPwshHookQuotesSelfPathWithSpaces(t *testing.T) {
	selfPath := "C:/Users/Jane Doe/direnv/direnv.exe"

	hookStr, err := Pwsh.Hook()
	if err != nil {
		t.Fatalf("Hook() returned error: %v", err)
	}

	hookTemplate, err := template.New("hook").Parse(hookStr)
	if err != nil {
		t.Fatalf("failed to parse hook template: %v", err)
	}

	var out strings.Builder
	err = hookTemplate.Execute(&out, HookContext{SelfPath: selfPath})
	if err != nil {
		t.Fatalf("failed to execute hook template: %v", err)
	}

	expected := `$export = (& "C:/Users/Jane Doe/direnv/direnv.exe" export pwsh) -join [Environment]::NewLine;`
	if !strings.Contains(out.String(), expected) {
		t.Errorf("expected hook output to contain %q, got:\n%s", expected, out.String())
	}
}
