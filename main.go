package main

import (
	"os"
)

// Configured at compile time
var bashPath string

func main() {
	// We drop $PWD from caller since it can include symlinks, which will
	// break relative path access when finding .envrc in a parent.
	_ = os.Unsetenv("PWD")

	var env = GetEnv()
	var args = os.Args

	setupLogging(env)

	err := CommandsDispatch(env, args)
	if err != nil {
		logError("error %v", err)
		os.Exit(1)
	}
}
