package main

import (
	"fmt"
)

var CmdHelp = &Cmd{
	Name: "help",
	Desc: "shows this help",
	Fn: func(env Env, args []string) (err error) {
		fmt.Printf(`direnv v%s
Usage: direnv COMMAND [...ARGS]

Available commands
------------------
`, VERSION)
		for _, cmd := range CmdList {
			if !cmd.Private {
				fmt.Printf("%s:\n  %s\n", cmd.Name, cmd.Desc)
			}
		}
		return
	},
}
