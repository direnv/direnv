package main

import (
	"fmt"
)

var CmdHelp = &Cmd{
	Name: "help",
	Desc: "shows this help",
	Fn: func(env Env, args []string) (err error) {
		fmt.Println("Help!")
		return
	},
}
