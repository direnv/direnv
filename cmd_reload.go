package main

import (
	"fmt"
)

// CmdReload is `direnv reload`
var CmdReload = &Cmd{
	Name: "reload",
	Desc: "triggers an env reload",
	Action: actionWithConfig(func(env Env, args []string, config *Config) error {
		foundRC, err := config.FindRC()
		if err != nil {
			return err
		}
		if foundRC == nil {
			return fmt.Errorf(".envrc not found")
		}

		return foundRC.Touch()
	}),
}
