package main

import (
	"path/filepath"
)

// Loosely
// http://standards.freedesktop.org/basedir-spec/basedir-spec-0.8.html
// We don't handle XDG_CONFIG_DIRS yet
func XdgConfigDir(env Env, programName string) string {
	if env["XDG_CONFIG_HOME"] != "" {
		return filepath.Join(env["XDG_CONFIG_HOME"], programName)
	} else if env["HOME"] != "" {
		return filepath.Join(env["HOME"], ".config", programName)
	}
	// In theory we could also read /etc/passwd and look for the home matching the process' Uid
	return ""
}
