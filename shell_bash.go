package main

type BASH struct{}

func (b *BASH) Hook() string {
	return `PROMPT_COMMAND="eval \$(direnv private export bash);$PROMPT_COMMAND"`
}

func (b *BASH) Escape(str string) string {
	return ShellEscape(str)
}
