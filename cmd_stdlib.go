package main

import (
	"fmt"
)

// `direnv stdlib`
var CmdStdlib = &Cmd{
	Name: "stdlib",
	Desc: "Displays the stdlib available in the .envrc execution context",
	Fn: func(env Env, args []string) (err error) {
		var config *Config
		if config, err = LoadConfig(env); err != nil {
			return
		}

		fmt.Printf(STDLIB, config.SelfPath)
		return
	},
}
