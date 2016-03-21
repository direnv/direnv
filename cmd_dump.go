package main

import "fmt"

// `direnv dump`
var CmdDump = &Cmd{
	Name:    "dump",
	Desc:    "Used to export the inner bash state at the end of execution",
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		fmt.Println(env.Serialize())
		return
	},
}
