package main

import (
	"fmt"
)

var CmdReload = &Cmd{
	Name: "reload",
	Desc: "triggers an env reload",
	Fn: func(env Env, args []string) error {
		config, err := LoadConfig(env)
		if err != nil {
			return err
		}

		foundRC := config.FindRC()
		if foundRC == nil {
			return fmt.Errorf(".envrc not found")
		}

		return foundRC.Touch()
	},
}
