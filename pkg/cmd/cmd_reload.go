package cmd

import (
	"fmt"
)

// CmdReload is `direnv reload`
var CmdReload = &Cmd{
	Name: "reload",
	Desc: "Triggers an env reload",
	Action: actionWithConfig(func(_ Env, _ []string, config *Config) error {
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
