package main

import (
	"fmt"
	"os"
)

// CmdAllow is `direnv allow [PATH_TO_RC]`
var CmdAllow = &Cmd{
	Name:   "allow",
	Desc:   "Grants direnv to load the given .envrc",
	Args:   []string{"[PATH_TO_RC]"},
	Action: actionWithConfig(cmdAllowAction),
}

func cmdAllowAction(env Env, args []string, config *Config) (err error) {
	var rcPath string
	if len(args) > 1 {
		rcPath = args[1]
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	rc := FindRC(rcPath, config)
	if rc == nil {
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Allow()
}
