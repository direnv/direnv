package main

import (
	"fmt"
	"os"
)

func main() {
	env := GetEnv()

	if err := CommandsDispatch(env, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "direnv %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}
}
