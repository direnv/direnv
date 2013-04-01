package main

// ZSH is a singleton instance of ZSH_T
type zsh int

var ZSH zsh

func (z zsh) Hook() string {
	return `
direnv_hook() { eval "$(direnv export zsh)"; };
[[ -z $precmd_functions ]] && precmd_functions=();
precmd_functions=($precmd_functions direnv_hook)
	`
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
