package cmd

import (
	"errors"
)

// systemdShell is not a real shell
type systemdShell struct{}

// Systemd is not really a shell but is useful to add support
// to systemd EnvironmentFile(https://0pointer.de/public/systemd-man/systemd.exec.html#EnvironmentFile=)
var Systemd Shell = systemdShell{}

func (sh systemdShell) Hook() (string, error) {
	return "", errors.New("this feature is not supported")
}

func (sh systemdShell) Export(e ShellExport) (out string) {
	for key, value := range e {
		if value != nil {
			out += sh.export(key, *value)
		}
	}
	return out
}

func (sh systemdShell) Dump(env Env) (out string) {
	for key, value := range env {
		out += sh.export(key, value)
	}
	return out
}

func (sh systemdShell) export(key, value string) string {
	return key + "=" + "\"" + value + "\"\n"
}
