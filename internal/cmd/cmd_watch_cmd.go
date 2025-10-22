package cmd

import (
	"fmt"
	"os"
)

// CmdWatchCmd is `direnv watch-cmd SHELL CMD`
var CmdWatchCmd = &Cmd{
	Name:    "watch-cmd",
	Desc:    "Watch a command output for changes",
	Args:    []string{"SHELL", "CMD"},
	Private: true,
	Action:  actionWithConfig(watchCmdCommand),
}

func watchCmdCommand(env Env, args []string, config *Config) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("a command is required to add to the list of watches")
	}

	shellName := args[1]

	shell := DetectShell(shellName)

	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", shellName)
	}

	watches := NewCmdValues()
	watchString, ok := env[DIRENV_CMD_WATCHES]
	if ok {
		err = watches.Unmarshal(watchString)
		if err != nil {
			return err
		}
	}

	err = watches.Update(args[2], config)
	if err != nil {
		return
	}

	e := make(ShellExport)
	e.Add(DIRENV_CMD_WATCHES, watches.Marshal())

	os.Stdout.WriteString(shell.Export(e))

	return
}

