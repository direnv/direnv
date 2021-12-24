package cmd

import (
	"os"
)

var (
	bashPath string
	stdlib   string
	version  string
)

// Main is the main entrypoint to direnv
func Main(env Env, args []string, modBashPath string, modStdlib string, modVersion string) error {
	// We drop $PWD from caller since it can include symlinks, which will
	// break relative path access when finding .envrc or .env in a parent.
	_ = os.Unsetenv("PWD")

	setupLogging(env)
	bashPath = modBashPath
	stdlib = modStdlib
	version = modVersion

	err := CommandsDispatch(env, args)
	if err != nil {
		logError("error %v", err)
	}
	return err
}
