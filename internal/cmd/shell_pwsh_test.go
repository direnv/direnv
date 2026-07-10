package cmd

import (
	"strings"
	"testing"
	"text/template"
)

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
