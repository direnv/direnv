package cmd

import (
	"strings"

	"github.com/direnv/direnv/v2/gzenv"
)

// IgnoredKeys is list of keys we don't want to deal with
var IgnoredKeys = map[string]bool{
	// direnv env config
	"DIRENV_CONFIG": true,
	"DIRENV_BASH":   true,

	// should only be available inside of the .envrc or .env
	"DIRENV_IN_ENVRC": true,

	"COMP_WORDBREAKS": true, // Avoids segfaults in bash
	"PS1":             true, // PS1 should not be exported, fixes problem in bash

	// variables that should change freely
	"OLDPWD":    true,
	"PWD":       true,
	"SHELL":     true,
	"SHELLOPTS": true,
	"SHLVL":     true,
	"_":         true,
}

// EnvDiff represents the diff between two environments
type EnvDiff struct {
	Prev map[string]string `json:"p"`
	Next map[string]string `json:"n"`
}

// NewEnvDiff is an empty constructor for EnvDiff
func NewEnvDiff() *EnvDiff {
	return &EnvDiff{make(map[string]string), make(map[string]string)}
}

// BuildEnvDiff analyses the changes between 'e1' and 'e2' and builds an
// EnvDiff out of it.
func BuildEnvDiff(e1, e2 Env) *EnvDiff {
	diff := NewEnvDiff()

	in := func(key string, e Env) bool {
		_, ok := e[key]
		return ok
	}

	for key := range e1 {
		if IgnoredEnv(key) {
			continue
		}
		if e2[key] != e1[key] || !in(key, e2) {
			diff.Prev[key] = e1[key]
		}
	}

	for key := range e2 {
		if IgnoredEnv(key) {
			continue
		}
		if e2[key] != e1[key] || !in(key, e1) {
			diff.Next[key] = e2[key]
		}
	}

	return diff
}

// LoadEnvDiff unmarshalls a gzenv string back into an EnvDiff.
func LoadEnvDiff(gzenvStr string) (diff *EnvDiff, err error) {
	diff = new(EnvDiff)
	err = gzenv.Unmarshal(gzenvStr, diff)
	return
}

// Any returns if the diff contains any changes.
func (diff *EnvDiff) Any() bool {
	return len(diff.Prev) > 0 || len(diff.Next) > 0
}

// ToShell applies the env diff as a set of commands that are understood by
// the target `shell`. The outputted string is then meant to be evaluated in
// the target shell.
func (diff *EnvDiff) ToShell(shell Shell) string {
	e := make(ShellExport)

	for key := range diff.Prev {
		_, ok := diff.Next[key]
		if !ok {
			e.Remove(key)
		}
	}

	for key, value := range diff.Next {
		e.Add(key, value)
	}

	return shell.Export(e)
}

// Patch applies the diff to the given env and returns a new env with the
// changes applied.
func (diff *EnvDiff) Patch(env Env) (newEnv Env) {
	newEnv = make(Env)

	for k, v := range env {
		newEnv[k] = v
	}

	for key := range diff.Prev {
		delete(newEnv, key)
	}

	for key, value := range diff.Next {
		newEnv[key] = value
	}

	return newEnv
}

// Reverse flips the diff so that it applies the other way around.
func (diff *EnvDiff) Reverse() *EnvDiff {
	return &EnvDiff{diff.Next, diff.Prev}
}

// Serialize marshalls the environment diff to the gzenv format.
func (diff *EnvDiff) Serialize() string {
	return gzenv.Marshal(diff)
}

//// Utils

// IgnoredEnv returns true if the key should be ignored in environment diffs.
func IgnoredEnv(key string) bool {
	if strings.HasPrefix(key, "__fish") {
		return true
	}
	if strings.HasPrefix(key, "BASH_FUNC_") {
		return true
	}
	_, found := IgnoredKeys[key]
	return found
}
