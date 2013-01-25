package env

import (
	"github.com/zimbatm/direnv/src/shell"
)

type EnvDiff map[string]string

func ToBash(env EnvDiff) string {
	str := ""
	for key, value := range env {
		// FIXME: This is not exacly as the ruby nil
		if value == "" {
			if key == "PS1" {
				// unsetting PS1 doesn't restore the default in OSX's bash
			} else {
				str += "unset " + key + ";"
			}
		} else {
			str += "export " + key + "=" + shell.Escape(value) + ";"
		}
	}
	return str
}

func Diff(env1 map[string]string, env2 map[string]string) EnvDiff {
	envDiff := make(EnvDiff)

	for key, _ := range env1 {
		if env2[key] != env1[key] && !IgnoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	// FIXME: I'm sure there is a smarter way to do that
	for key, _ := range env2 {
		if env2[key] != env1[key] && !IgnoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	return envDiff
}
