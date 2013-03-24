package main

type ZSH struct{}

func (z *ZSH) Hook() string {
	return `
	direnv_hook() { eval \$(direnv private export zsh) };
	[[ -z $precmd_functions ]] && precmd_functions=();
	precmd_functions=($precmd_functions direnv_hook)
	`
}

func (z *ZSH) Escape(str string) string {
	return ShellEscape(str)
}
