package main

import (
	"fmt"
	"strings"
)

// CmdStdlib is `direnv stdlib`
var CmdStdlib = &Cmd{
	Name: "stdlib",
	Desc: "Displays the stdlib available in the .envrc execution context",
	Action: actionWithConfig(func(env Env, args []string, config *Config) error {
		fmt.Println(strings.Replace(StdLib, "$(command -v direnv)", config.SelfPath, 1))
		return nil
	}),
}
