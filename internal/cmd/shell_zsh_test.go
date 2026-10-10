package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func TestZshHookEvalContext(t *testing.T) {
	shellPath, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("zsh is not installed")
	}
	tmpl, err := template.New("hook").Parse(zshHook)
	if err != nil {
		t.Fatal(err)
	}
	var hook bytes.Buffer
	if err := tmpl.Execute(&hook, HookContext{SelfPath: "__direnv_export"}); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(t.TempDir(), "calls")

	// zsh runs precmd between the lines it reads, so each check sits on one line.
	script := `__direnv_export() { print >> ` + log + `; }
` + hook.String() + `
: > ` + log + `
[[ -s ` + log + ` ]] || { echo "hook skipped in precmd"; exit 1; }
: > ` + log + `; ( cd / ); x=$(cd /); [[ ! -s ` + log + ` ]] || { echo "hook ran in a subshell"; exit 1; }
exit 0
`
	command := exec.Command(shellPath, "-f", "-i")
	command.Stdin = strings.NewReader(script)
	command.Env = append(os.Environ(), "ZDOTDIR="+t.TempDir())
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, output)
	}
}
