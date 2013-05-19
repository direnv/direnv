package main

import (
	"fmt"
	"os"
)

func main() {
	var env = GetEnv()
	var args = os.Args

	err := CommandsDispatch(env, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
