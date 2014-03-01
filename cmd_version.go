package main

import (
	"fmt"
)

var CmdVersion = &Cmd{
	Name:    "version",
	Desc:    "Prints the current direnv version",
	Private: true,
	Fn: func(env Env, args []string) error {
		fmt.Println(VERSION)
		return nil
	},
}
