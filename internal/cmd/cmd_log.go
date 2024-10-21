package cmd

import "fmt"

// CmdFetchURL is `direnv fetchurl <url> [<integrity-hash>]`
var CmdLog = &Cmd{
	Name:   "log",
	Desc:   "Logs a given message using log-related environment variables: [DIRENV_LOG_FORMAT, DIRENV_LOG_FILTER]",
	Args:   []string{"<message>"},
	Action: actionWithConfig(cmdLog),
}

func cmdLog(env Env, args []string, config *Config) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("missing message argument")
	}
	message := args[1]
	logStatus(env, message)
	return nil
}
