package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func fixMsys2Paths(env Env) {
	if env["MSYSTEM"] == "" {
		return
	}
	shell := strings.TrimSuffix(strings.ToLower(filepath.Base(env["SHELL"])), ".exe")
	switch shell {
	case "bash", "elvish", "fish", "murex", "tcsh", "zsh":
	default:
		return
	}
	path := env["PATH"]
	if !strings.Contains(path, ";") {
		return
	}
	cygpath, err := exec.LookPath("cygpath")
	if err != nil {
		panic(fmt.Sprintf("direnv: MSYS2 detected but cygpath not found: %v", err))
	}
	out, err := exec.Command(cygpath, "--path", "--unix", "--", path).Output()
	if err != nil {
		panic(fmt.Sprintf("direnv: cygpath failed: %v", err))
	}
	env["PATH"] = strings.TrimSpace(string(out))
}
