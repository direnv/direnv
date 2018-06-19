package main

import (
	"fmt"
)

var CmdVersion = &Cmd{
	Name: "version",
	Desc: "prints the version",
	Fn: func(env Env, args []string) error {
		fmt.Println(VERSION)
		return nil
	},
}
