package main

import (
	"fmt"
)

// CmdVersion is `direnv version`
var CmdVersion = &Cmd{
	Name:    "version",
	Desc:    "prints the version (" + Version + ")",
	Aliases: []string{"--version"},
	Action: actionSimple(func(env Env, args []string) error {
		fmt.Println(Version)
		return nil
	}),
}
