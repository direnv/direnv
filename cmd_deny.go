package main

import (
	"fmt"
	"os"
)

// `direnv deny [path_to_rc]`
var CmdDeny = &Cmd{
	Name: "deny",
	Desc: "Revokes the auhorization of a given .envrc",
	Args: []string{"[path_to_rc]"},
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
