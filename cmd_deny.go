package main

import (
	"fmt"
	"os"
)

// CmdDeny is `direnv deny [PATH_TO_RC]`
var CmdDeny = &Cmd{
	Name:   "deny",
	Desc:   "Revokes the authorization of a given .envrc",
	Args:   []string{"[PATH_TO_RC]"},
	Action: actionWithConfig(cmdDenyAction),
}

func cmdDenyAction(env Env, args []string, config *Config) (err error) {
	var rcPath string

	if len(args) > 1 {
		rcPath = args[1]
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	rc, err := FindRC(rcPath, config)
	if err != nil {
		return err
	} else if rc == nil {
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Deny()
}
