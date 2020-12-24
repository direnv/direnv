// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris js,wasm

package lpenv

import "strings"

// Getenv looks for the key in the given environment variables
func Getenv(env []string, key string) string {
	if len(key) == 0 {
		return ""
	}

	prefix := key + "="
	for _, pair := range env {
		if strings.HasPrefix(pair, prefix) {
			return pair[len(prefix):]
		}

	}
	return ""
}
