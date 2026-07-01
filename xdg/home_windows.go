package xdg

func homeDir(env map[string]string) string {
	if env["HOME"] != "" {
		return env["HOME"]
	}
	return env["USERPROFILE"]
}
