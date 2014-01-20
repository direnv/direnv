package main

type bash int

var BASH bash

const BASH_HOOK = `
_direnv_hook() {
  eval "$(direnv export bash)";
};
if ! [[ "$PROMPT_COMMAND" =~ _direnv_hook ]]; then
  PROMPT_COMMAND="_direnv_hook;$PROMPT_COMMAND";
fi
`

func (b bash) Hook() string {
	return BASH_HOOK
}

func (b bash) Escape(str string) string {
	return ShellEscape(str)
}

func (b bash) Export(key, value string) string {
	return "export " + b.Escape(key) + "=" + b.Escape(value) + ";"
}

func (b bash) Unset(key string) string {
	return "unset " + b.Escape(key) + ";"
}
