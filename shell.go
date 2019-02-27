package main

import (
	"path/filepath"
)

/*
 * Shells
 */
type Shell interface {
	Hook() (string, error)
	Export(e ShellExport, q ShellQuotes) string
	Dump(env Env) string
}

// Used to describe what to generate for the shell
type ShellExport map[string]*string

func (e ShellExport) Add(key, value string) {
	e[key] = &value
}

func (e ShellExport) Remove(key string) {
	e[key] = nil
}

// Shell-specific commands
type ShellQuotes map[Shell][]string

func MergeShellQuotes(l, r ShellQuotes) (merged ShellQuotes) {
	merged = make(ShellQuotes)
	for k, v := range l {
		merged[k] = v
	}
	for k, v := range r {
		merged[k] = append(merged[k], v...)
	}
	return
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
	case "gzenv":
		return GZENV
	case "vim":
		return VIM
	case "tcsh":
		return TCSH
	case "json":
		return JSON
	case "elvish":
		return ELVISH
	}

	return nil
}
