package cmd

import (
	"fmt"
	"strings"
)

type gha struct{}

// GitHubActions shell instance
var GitHubActions Shell = gha{}

func (sh gha) Hook() (string, error) {
	return "", fmt.Errorf("Hook not implemented for GitHub Actions shell")
}

func (sh gha) Export(e ShellExport) string {
	var b strings.Builder
	for key, value := range e {
		if value == nil {
			sh.unset(&b, key)
		} else {
			sh.export(&b, key, *value)
		}
	}
	return b.String()
}

const ghaDelim = "DIRENV_GITHUB_ACTIONS_EOV\n"

func (sh gha) Dump(env Env) string {
	var b strings.Builder

	for key, value := range env {
		sh.export(&b, key, value)
	}
	return b.String()
}

func (sh gha) export(b *strings.Builder, key, value string) {
	b.WriteString(key)
	b.WriteString("<<")
	b.WriteString(ghaDelim)
	b.WriteString(value)
	if value != "" && !strings.HasSuffix(value, "\n") {
		b.WriteByte('\n')
	}
	b.WriteString(ghaDelim)
}

func (sh gha) unset(b *strings.Builder, key string) {
	sh.export(b, key, "")
}
