package cmd

import "strings"

// ZSH is a singleton instance of ZSH_T
type zsh struct{}

// Zsh adds support for the venerable Z shell.
var Zsh Shell = zsh{}

// The guard skips the hook when chpwd fires within a non-interactive
// shell or a subshell, e.g., from completion. It is inspired by zsh's
// Functions/Chpwd/chpwd_recent_dirs but does not skip a non-toplevel
// ZSH_EVAL_CONTEXT. See https://github.com/direnv/direnv/issues/1633 as to why.
const zshHook = `
_direnv_hook() {
  setopt localoptions localtraps
  if [[ ! -o interactive || $ZSH_SUBSHELL -ne 0 ]]; then
    return
  fi
  vars="$({{.BashSelfPath}} export zsh)"
  trap -- '' SIGINT
  eval "$vars"
  trap - SIGINT
}
typeset -ag precmd_functions
if (( ! ${precmd_functions[(I)_direnv_hook]} )); then
  precmd_functions=(_direnv_hook $precmd_functions)
fi
typeset -ag chpwd_functions
if (( ! ${chpwd_functions[(I)_direnv_hook]} )); then
  chpwd_functions=(_direnv_hook $chpwd_functions)
fi
`

func (sh zsh) Hook() (string, error) {
	return zshHook, nil
}

func (sh zsh) Export(e ShellExport) (string, error) {
	var out strings.Builder
	for key, value := range e {
		if value == nil {
			out.WriteString(sh.unset(key))
		} else {
			out.WriteString(sh.export(key, *value))
		}
	}
	return out.String(), nil
}

func (sh zsh) Dump(env Env) (string, error) {
	var out strings.Builder
	for key, value := range env {
		out.WriteString(sh.export(key, value))
	}
	return out.String(), nil
}

func (sh zsh) export(key, value string) string {
	return "export " + sh.escape(key) + "=" + sh.escape(value) + ";"
}

func (sh zsh) unset(key string) string {
	return "unset " + sh.escape(key) + ";"
}

func (sh zsh) escape(str string) string {
	return BashEscape(str)
}
