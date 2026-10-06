// Package xdg is a minimal implementation of the XDG specification.
//
// https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
package xdg

import (
	"path/filepath"
	"strings"
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

// ConfigDirs returns all config directories for the application in priority
// order per the XDG Base Directory Specification. Returns user config dir
// first (XDG_CONFIG_HOME or ~/.config fallback), then system config dirs
// (XDG_CONFIG_DIRS or /etc/xdg fallback).
func ConfigDirs(env map[string]string, programName string) []string {
	var dirs []string

	if env["XDG_CONFIG_HOME"] != "" {
		dirs = append(dirs, filepath.Join(env["XDG_CONFIG_HOME"], programName))
	} else if env["HOME"] != "" {
		dirs = append(dirs, filepath.Join(env["HOME"], ".config", programName))
	}

	if env["XDG_CONFIG_DIRS"] != "" {
		for _, dir := range strings.Split(env["XDG_CONFIG_DIRS"], ":") {
			if dir != "" {
				dirs = append(dirs, filepath.Join(dir, programName))
			}
		}
	} else {
		dirs = append(dirs, filepath.Join("/etc/xdg", programName))
	}

	return dirs
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
