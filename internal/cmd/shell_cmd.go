package cmd

import (
	"fmt"
)

type windowscmd struct{}

// windowscmd shell instance
var WindowsCmd Shell = windowscmd{}

func (sh windowscmd) Hook() (string, error) {
	const hook = `@rem Use this 'denv' command to update environment variables for the current directory.
@rem direnv hook cmd> %USERPROFILE%\.local\bin\denv.bat
@
@echo @echo off> %TEMP%\%USERNAME%_vars.bat
@direnv export cmd>> %TEMP%\%USERNAME%_vars.bat
@%TEMP%\%USERNAME%_vars.bat
@del %TEMP%\%USERNAME%_vars.bat
`
	return hook, nil
}

func (sh windowscmd) Export(e ShellExport) (string, error) {
	var out string
	for key, value := range e {
		if key != "" {
			if value == nil {
				out += sh.unset(key)
			} else {
				out += sh.export(key, *value)
			}
		}
	}
	return out, nil
}

func (sh windowscmd) Dump(env Env) (string, error) {
	var out string
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out, nil
}

func (sh windowscmd) WindowsNative() bool {
	return true
}

func (sh windowscmd) export(key, value string) string {
	return fmt.Sprintf("set \"%s=%s\"\n", sh.escapeEnvKey(key), sh.escapeVerbatimString(value))
}

func (sh windowscmd) unset(key string) string {
	return fmt.Sprintf("set %s=\n", sh.escapeVerbatimEnvKey(key))
}

func (sh windowscmd) escapeEnvKey(str string) string {
	return PowerShellEscapeEnvKey(str)
}

func (windowscmd) escapeVerbatimEnvKey(str string) string {
	return PowerShellEscapeVerbatimEnvKey(str)
}

func (windowscmd) escapeVerbatimString(str string) string {
	return PowerShellEscapeVerbatimString(str)
}
