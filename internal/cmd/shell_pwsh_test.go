package cmd

import (
	"strings"
	"testing"
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
