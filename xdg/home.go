//go:build !windows

package xdg

func homeDir(env map[string]string) string {
	return env["HOME"]
}
