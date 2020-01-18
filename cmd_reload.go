package main

import (
	"fmt"
)

// CmdReload is `direnv reload`
var CmdReload = &Cmd{
	Name: "reload",
	Desc: "triggers an env reload",
	Action: actionWithConfig(func(env Env, args []string, config *Config) error {
		foundRC := config.FindRC()
		if foundRC == nil {
			return fmt.Errorf(".envrc not found")
		}

		return foundRC.Touch()
	}),
}
