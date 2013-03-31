package main

import (
	"fmt"
	"os"
)

func main() {
	env := GetEnv()

	if err := publicCommands(env, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
