package main

import (
	"fmt"
)

// `direnv hook $0`
var CmdHook = &Cmd{
	Name: "hook",
	Desc: "Used to setup the shell hook",
	Args: []string{"SHELL"},
	Fn: func(env Env, args []string) (err error) {
		var target string

		if len(args) > 1 {
			target = args[1]
		}

		shell := DetectShell(target)
		if shell == nil {
			return fmt.Errorf("Unknown target shell '%s'", target)
		}

		h, err := shell.Hook()
		if err != nil {
			return err
		}

		fmt.Println(h)

		return
	},
}
