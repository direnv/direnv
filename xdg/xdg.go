// Package xdg is a minimal implementation of the XDG specification.
//
// https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
package xdg

import (
	"path/filepath"
)

// DataDir returns the data folder for the application
func DataDir(env map[string]string, programName string) string {
	if env["XDG_DATA_HOME"] != "" {
		return filepath.Join(env["XDG_DATA_HOME"], programName)
	} else if env["HOME"] != "" {
		return filepath.Join(env["HOME"], ".local", "share", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}

// ConfigDir returns the config folder for the application
//
// The XDG_CONFIG_DIRS case is not being handled
func ConfigDir(env map[string]string, programName string) string {
	if env["XDG_CONFIG_HOME"] != "" {
		return filepath.Join(env["XDG_CONFIG_HOME"], programName)
	} else if env["HOME"] != "" {
		return filepath.Join(env["HOME"], ".config", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}

// CacheDir returns the cache directory for the application
func CacheDir(env map[string]string, programName string) string {
	if env["XDG_CACHE_HOME"] != "" {
		return filepath.Join(env["XDG_CACHE_HOME"], programName)
	} else if env["HOME"] != "" {
		return filepath.Join(env["HOME"], ".cache", programName)
	}
	// In theory we could also read /etc/passwd and look for the home based on
	// the process' UID
	return ""
}
