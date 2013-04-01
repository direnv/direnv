package main

import (
	"fmt"
)

var CmdVersion = &Cmd{
	Name:    "version",
	Private: true,
	Fn: func(env Env, args []string) error {
		fmt.Println(VERSION)
		return nil
	},
}
