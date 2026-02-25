package cmd

import (
	"regexp"
	"strings"
)

// Hooks contains all hooks in the form: map[shellName -> map[hookName -> hook]]
type Hooks map[string]map[string]string

var (
	shellEscape = "direnv-shell-escape"
	shellEscapePrefix = "{{{" + shellEscape + " "
	shellEscapeSuffix = " " + shellEscape + "}}}"
	shellEscapeRegex = regexp.MustCompile(`(?s)` + shellEscapePrefix + `.*?` + shellEscapeSuffix)
)

// Get gets a hook
func (hooks Hooks) Get(hookName string, hookableShell HookableShell) string {
	hooksForShell := hooks[hookableShell.Name()]
	if hooksForShell != nil {
		// PERF: Rather than process the escapes in all the hooks ahead of time,
		// we process them lazily here
		return processEscapes(hooksForShell[hookName], hookableShell)
	}

	return ""
}

// Set sets a hook
func (hooks Hooks) Set(hookName string, hook string, shellName string) {
	hooksForShell := hooks[shellName]
	if hooksForShell == nil {
		hooksForShell = map[string]string{}
		hooks[shellName] = hooksForShell
	}

	hooksForShell[hookName] = hook
}

func processEscapes(hook string, hookableShell HookableShell) string {
	return shellEscapeRegex.ReplaceAllStringFunc(hook, func(match string) string {
		matchWithoutPrefix := strings.TrimPrefix(match, shellEscapePrefix)
		matchWithoutPrefixAndSuffix := strings.TrimSuffix(matchWithoutPrefix, shellEscapeSuffix)
		return hookableShell.Escape(matchWithoutPrefixAndSuffix)
	})
}
