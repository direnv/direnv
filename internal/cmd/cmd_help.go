package cmd

import (
	"fmt"
	"strings"
)

// CmdHelp is `direnv help`
var CmdHelp = &Cmd{
	Name:    "help",
	Desc:    "Shows this help",
	Args:    []string{"[SHOW_PRIVATE]"},
	Aliases: []string{"--help"},
	Action: actionSimple(func(_ Env, args []string) (err error) {
		var showPrivate = len(args) > 1
		fmt.Printf(`direnv v%s
Usage: direnv COMMAND [...ARGS]

Available commands
------------------
`, version)
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
				fmt.Printf("%s%s:\n", cmd.Name, opts)
				for _, alias := range cmd.Aliases {
					if alias[0:1] != "-" {
						fmt.Printf("%s%s:\n", alias, opts)
					}
				}
				fmt.Printf("  %s\n", cmd.Desc)
			}
		}

		if showPrivate {
			fmt.Println("* = private commands")
		}
		return
	}),
}
