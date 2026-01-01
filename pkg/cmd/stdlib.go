package cmd

import "strings"

// getStdlib returns the stdlib.sh, with references to direnv replaced.
func getStdlib(config *Config) string {
	return strings.Replace(stdlib, "$(command -v direnv)", config.SelfPath, 1)
}
