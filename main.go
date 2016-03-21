package main

import (
	"os"
)

func main() {
	var env = GetEnv()
	var args = os.Args

	setupLogging(env)

	err := CommandsDispatch(env, args)
	if err != nil {
		log_error("error %v", err)
		os.Exit(1)
	}
}
