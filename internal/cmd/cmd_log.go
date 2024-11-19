package cmd

import (
	"errors"
	"fmt"
)

// CmdLog is `direnv log [--status | --error] <message>`
var CmdLog = &Cmd{
	Name:   "log",
	Desc:   "Logs a given message",
	Args:   []string{"[--status | --error]", "<message>"},
	Action: actionWithConfig(cmdLog),
}

func cmdLog(_ Env, args []string, c *Config) (err error) {
	if len(args) != 3 {
		return errors.New("invalid arguments")
	}
	logType := args[1]
	message := args[2]
	if logType == "--status" || logType == "-status" {
		logStatus(c, message)
	} else if logType == "--error" || logType == "-error" {
		logError(c, message)
	} else {
		return fmt.Errorf("invalid log-type '%s'", logType)
	}
	return nil
}
