package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"text/template"
)

const weirdDirName = "a b $y `x` 'q'"

func TestStdlibEscapesSelfPath(t *testing.T) {
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is not installed")
	}
	saved := stdlib
	t.Cleanup(func() { stdlib = saved })
	stdlib = `direnv="$(command -v direnv)"`

	selfPath := filepath.Join(t.TempDir(), weirdDirName, "direnv")
	cmd := exec.Command(bashPath, "--noprofile", "--norc", "-s")
	cmd.Stdin = strings.NewReader(getStdlib(&Config{SelfPath: selfPath}) + "\nprintf %s \"$direnv\"")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != selfPath {
		t.Fatalf("expected %q, got %q", selfPath, out)
	}
}

func TestBashHookEscapesSelfPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires a unix shell script")
	}
	bashPath, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is not installed")
	}
	dir := filepath.Join(t.TempDir(), weirdDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	selfPath := filepath.Join(dir, "direnv")
	if err := os.WriteFile(selfPath, []byte("#!/bin/sh\necho export LOADED=yes\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	hook, err := Bash.Hook()
	if err != nil {
		t.Fatal(err)
	}
	var rendered bytes.Buffer
	if err := template.Must(template.New("hook").Parse(hook)).Execute(&rendered, HookContext{SelfPath: selfPath}); err != nil {
		t.Fatal(err)
	}
	script := rendered.String() + "\n_direnv_hook\nprintf %s \"$LOADED\""
	out, err := exec.Command(bashPath, "--noprofile", "--norc", "-c", script).CombinedOutput()
	if err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if string(out) != "yes" {
		t.Fatalf("expected the hook to run %q, got %q", selfPath, out)
	}
}
