package cmd

import "strings"

// getStdlib returns the stdlib.sh, with references to direnv replaced.
func getStdlib(config *Config) string {
	str := strings.Replace(stdlib, "${DIRENV_SYSCONFIG:-/etc/direnv}", "${DIRENV_SYSCONFIG:-"+config.SysConfDir+"}", 1)

	return strings.Replace(str, "$(command -v direnv)", config.SelfPath, 1)
}
