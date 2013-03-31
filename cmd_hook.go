package main

import (
	"fmt"
)

// `direnv hook $0`
// $0 starts with "-" and go tries to parse it as an argument
//
// This command is public for historical reasons
func Hook(env Env, args []string) (err error) {
	var target string

	if len(args) > 1 {
		target = args[1]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("Unknown target shell '%s'", target)
	}

	fmt.Println(shell.Hook())

	return
}
