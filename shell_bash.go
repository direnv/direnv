package main

type bash int

var BASH bash

func (b bash) Hook() string {
	return `PROMPT_COMMAND="eval \"\$(direnv export bash) ;$PROMPT_COMMAND\""`
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
