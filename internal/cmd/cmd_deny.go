package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// CmdDeny is `direnv deny [PATH_TO_RC]`
var CmdDeny = &Cmd{
	Name:    "block",
	Desc:    "Revokes the authorization of a given .envrc or .env file.",
	Args:    []string{"[PATH_TO_RC]"},
	Aliases: []string{"deny", "disallow", "revoke"},
	Action:  actionWithConfig(cmdDenyAction),
}

func cmdDenyAction(_ Env, args []string, config *Config) (err error) {
	var rcPath string

	if len(args) > 1 {
		if rcPath, err = filepath.Abs(args[1]); err != nil {
			return err
		}
		if rcPath, err = filepath.EvalSymlinks(rcPath); err != nil {
			return err
		}
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	rc, err := FindRC(rcPath, config)
	if err != nil {
		return err
	} else if rc == nil {
		if config.LoadDotenv {
			return fmt.Errorf(".envrc or .env file not found")
		}
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Deny()
}
