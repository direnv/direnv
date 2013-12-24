package main

import (
	"fmt"
	"strings"
)

// `direnv help`
var CmdHelp = &Cmd{
	Name:    "help",
	Desc:    "shows this help",
	Args:    []string{"[SHOW_PRIVATE]"},
	Aliases: []string{"--help"},
	Fn: func(env Env, args []string) (err error) {
		var showPrivate = len(args) > 1
		fmt.Printf(`direnv v%s
Usage: direnv COMMAND [...ARGS]

Available commands
------------------
`, VERSION)
		for _, cmd := range CmdList {
			var opts string
			if len(cmd.Args) > 0 {
				opts = " " + strings.Join(cmd.Args, " ")
			}
			if cmd.Private {
				if showPrivate {
					fmt.Printf("*%s%s:\n  %s\n", cmd.Name, opts, cmd.Desc)
				}
			} else {
				fmt.Printf("%s%s:\n  %s\n", cmd.Name, opts, cmd.Desc)
			}
		}

		if showPrivate {
			fmt.Println("* = private commands")
		}
		return
	},
}
