package main

import "fmt"

// `direnv dump`
var CmdDump = &Cmd{
	Name:    "dump",
	Desc:    "Used to export the inner bash state at the end of execution",
	Args:    []string{"[SHELL]"},
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		target := "gzenv"

		if len(args) > 1 {
			target = args[1]
		}

		shell := DetectShell(target)
		if shell == nil {
			return fmt.Errorf("Unknown target shell '%s'", target)
		}

		fmt.Println(shell.Dump(env))

		return
	},
}
