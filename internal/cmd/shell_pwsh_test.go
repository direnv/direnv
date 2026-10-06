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

func TestPwshHookQuotesSelfPath(t *testing.T) {
	hook, err := Pwsh.Hook()
	if err != nil {
		t.Fatal(err)
	}
	tmpl, err := template.New("hook").Parse(hook)
	if err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := tmpl.Execute(&out, HookContext{SelfPath: "C:/Users/Jane O'Doe/$x/direnv.exe"}); err != nil {
		t.Fatal(err)
	}
	expected := `$export = (& 'C:/Users/Jane O''Doe/$x/direnv.exe' export pwsh) -join [Environment]::NewLine;`
	if !strings.Contains(out.String(), expected) {
		t.Errorf("expected hook output to contain %q, got:\n%s", expected, out.String())
	}
}
