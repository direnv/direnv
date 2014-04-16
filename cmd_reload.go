package main

import (
	"fmt"
)

var CmdReload = &Cmd{
	Name: "reload",
	Desc: "Triggers an env reload",
	Fn: func(env Env, args []string) error {
		config, err := LoadConfig(env)
		if err != nil {
			return err
		}

		foundRC := config.FindRC()
		if foundRC != nil {
			return foundRC.Touch()
		} else {
			return fmt.Errorf(".envrc not found")
		}

		return nil
	},
}
