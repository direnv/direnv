package main

import (
	"os"
)

func main() {
	var env = GetEnv()
	var args = os.Args

	err := CommandsDispatch(env, args)
	if err != nil {
		log("error %v", err)
		os.Exit(1)
	}
}
