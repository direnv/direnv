package cmd

import (
	"errors"
	"fmt"
)

// CmdLog is `direnv log [--status | --error] <message>`
var CmdLog = &Cmd{
	Name:   "log",
	Desc:   "Logs a given message using log-related environment variables: [DIRENV_LOG_FORMAT, DIRENV_LOG_FILTER]",
	Args:   []string{"[--status | --error]", "<message>"},
	Action: actionWithConfig(cmdLog),
}

func cmdLog(env Env, args []string, config *Config) (err error) {
	if len(args) != 3 {
		return errors.New("invalid arguments")
	}
	logType := args[1]
	message := args[2]
	if logType == "--status" || logType == "-status" {
		logStatus(env, message)
	} else if logType == "--error" || logType == "-error" {
		logError(message)
	} else {
		return fmt.Errorf("invalid log-type '%s'", logType)
	}
	return nil
}
