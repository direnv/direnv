package main

// ZSH is a singleton instance of ZSH_T
type zsh int

var ZSH zsh

const ZSH_HOOK = `
_direnv_hook() {
  eval "$(direnv export zsh)";
}
autoload -U add-zsh-hook
if [[ -z $precmd_functions[(r)_direnv_hook] ]]; then
  add-zsh-hook precmd _direnv_hook;
fi
`

func (z zsh) Hook() string {
	return ZSH_HOOK
}

func (z zsh) Escape(str string) string {
	return ShellEscape(str)
}

func (z zsh) Export(key, value string) string {
	return "export " + z.Escape(key) + "=" + z.Escape(value) + ";"
}

func (z zsh) Unset(key string) string {
	return "unset " + z.Escape(key) + ";"
}
