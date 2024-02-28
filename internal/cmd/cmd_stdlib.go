package cmd

import (
	"fmt"
)

// CmdStdlib is `direnv stdlib`
var CmdStdlib = &Cmd{
	Name: "stdlib",
	Desc: "Displays the stdlib available in the .envrc execution context",
	Action: actionWithConfig(func(_ Env, _ []string, config *Config) error {
		fmt.Println(getStdlib(config))
		return nil
	}),
}
