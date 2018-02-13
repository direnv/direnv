package main

import (
	"fmt"
	"os"
)

// `direnv whitelist [PATH_TO_RC]`
var CmdWhitelist = &Cmd{
	Name: "whitelist",
	Desc: "Permanently allows any .envrc that resides at the given path",
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

		rc := FindRC(rcPath, config.AllowDir(), config.WhitelistDir())
		if rc == nil {
			return fmt.Errorf(".envrc file not found")
		}
		return rc.Whitelist()
	},
}
