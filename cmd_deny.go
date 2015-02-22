package main

import (
	"fmt"
	"os"
)

// `direnv deny [PATH_TO_RC]`
var CmdDeny = &Cmd{
	Name: "deny",
	Desc: "Revokes the authorization of a given .envrc",
	Args: []string{"[PATH_TO_RC]"},
	Fn: func(env Env, args []string) (err error) {
		var rcPath string
		var config *Config

		if len(args) > 1 {
			rcPath = args[1]
		} else {
			if rcPath, err = os.Getwd(); err != nil {
				return
			}
		}

		if config, err = LoadConfig(env); err != nil {
			return
		}

		rc := FindRC(rcPath, config.AllowDir())
		if rc == nil {
			return fmt.Errorf(".envrc file not found")
		}
		return rc.Deny()
	},
}
