package main

import (
	"path/filepath"
)

/*
 * Shells
 */
type Shell interface {
	Hook() (string, error)
	Export(e ShellExport) string
}

// Used to describe what to generate for the shell
type ShellExport map[string]*string

func (e ShellExport) Add(key, value string) {
	e[key] = &value
}

func (e ShellExport) Remove(key string) {
	e[key] = nil
}

func DetectShell(target string) Shell {
	target = filepath.Base(target)
	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	switch target {
	case "bash":
		return BASH
	case "zsh":
		return ZSH
	case "fish":
		return FISH
	case "vim":
		return VIM
	case "tcsh":
		return TCSH
	case "json":
		return JSON
	}

	return nil
}
