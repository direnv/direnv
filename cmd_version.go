package main

import (
	"fmt"
)

var CmdVersion = &Cmd{
	Name:    "version",
	Desc:    "prints the version",
	Aliases: []string{"--version"},
	Action: actionSimple(func(env Env, args []string) error {
		fmt.Println(VERSION)
		return nil
	}),
}
