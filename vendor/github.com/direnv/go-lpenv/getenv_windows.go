package lpenv

import "strings"

// Getenv looks for the key in the given environment variables
//
// The Windows implementation is case-insensitive.
func Getenv(env []string, key string) string {
	if len(key) == 0 {
		return ""
	}

	prefix := strings.ToLower(key + "=")
	for _, pair := range env {
		if len(pair) > len(prefix) && prefix == strings.ToLower(pair[:len(prefix)]) {
			return pair[len(prefix):]
		}
	}
	return ""
}
