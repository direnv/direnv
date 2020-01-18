package main

import (
	"path/filepath"
)

// Shell is the interface that represents the interaction with the host shell.
type Shell interface {
	// Hook is the string that gets evaluated into the host shell config and
	// setups direnv as a prompt hook.
	Hook() (string, error)

	// Export outputs the ShellExport as an evaluatable string on the host shell
	Export(e ShellExport) string

	// Dump outputs and evaluatable string that sets the env in the host shell
	Dump(env Env) string
}

// ShellExport represents environment variables to add and remove on the host
// shell.
type ShellExport map[string]*string

// Add represents the additon of a new environment variable
func (e ShellExport) Add(key, value string) {
	e[key] = &value
}

// Remove represents the removal of a given `key` environment variable.
func (e ShellExport) Remove(key string) {
	e[key] = nil
}

// DetectShell returns a Shell instance from the given target.
//
// target is usually $0 and can also be prefixed by `-`
func DetectShell(target string) Shell {
	target = filepath.Base(target)
	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	switch target {
	case "bash":
		return Bash
	case "zsh":
		return Zsh
	case "fish":
		return Fish
	case "gzenv":
		return GzEnv
	case "vim":
		return Vim
	case "tcsh":
		return Tcsh
	case "json":
		return JSON
	case "elvish":
		return Elvish
	}

	return nil
}
