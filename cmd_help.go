package main

import (
	"fmt"
)

var CmdHelp = &Cmd{
	Name: "help",
	Desc: "shows this help",
	Fn: func(env Env, args []string) (err error) {
		fmt.Println("direnv COMMAND [...ARGS]")
		for _, cmd := range CmdList {
			if !cmd.Private {
				fmt.Printf("%s: %s\n", cmd.Name, cmd.Desc)
			}
		}
		return
	},
}
