package main

import (
	"fmt"
	"os"
)

func main() {
	env := GetEnv()

	if err := CommandsDispatch(env, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "direnv: %s\n", err)
		os.Exit(1)
	}
}
