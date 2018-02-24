package main

import (
	"fmt"
	"strings"
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

		fmt.Println(strings.Replace(STDLIB, "$(which direnv)", config.SelfPath, 1))
		return
	},
}
