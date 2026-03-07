package cmd

import (
	"path/filepath"
)

// Shell is the interface that represents the interaction with the host shell.
type Shell interface {
	Name() string

	// Hook is the string that gets evaluated into the host shell config and
	// setups direnv as a prompt hook.
	Hook() (string, error)

	// Export outputs the ShellExport as an evaluatable string on the host shell
	Export(e ShellExport) (string, error)

	// Dump outputs and evaluatable string that sets the env in the host shell
	Dump(env Env) (string, error)
}

// HookableShell is similar to the Shell interface, but with support for hooks
type HookableShell interface {
	Shell

	ExportWithHooks(shellExport ShellExport, hooks map[string]string, setProcessMarker *bool) (string, error)
}

// ShellExport represents environment variables to add and remove on the host
// shell.
type ShellExport map[string]*string

// Add represents the addition of a new environment variable
func (e ShellExport) Add(key, value string) {
	e[key] = &value
}

// Remove represents the removal of a given `key` environment variable.
func (e ShellExport) Remove(key string) {
	e[key] = nil
}

var supportedShellList = func() (shells map[string]Shell) {
	shells = map[string]Shell{}
	for _, shell := range []Shell{Bash, Elvish, Fish, GitHubActions, GzEnv, JSON, Murex, Tcsh, Vim, Zsh, Pwsh, Systemd} {
		shells[shell.Name()] = shell
	}

	return
}()

// DetectShell returns a Shell instance from the given target.
//
// target is usually $0 and can also be prefixed by `-`
func DetectShell(target string) Shell {
	target = filepath.Base(target)
	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	detectedShell, isValid := supportedShellList[target]
	if isValid {
		return detectedShell
	}
	return nil
}
