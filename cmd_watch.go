package main

import (
	"fmt"
	"os"
)

var CmdWatch = &Cmd{
	Name:    "watch",
	Desc:    "Adds a path to the list that direnv watches for changes",
	Args:    []string{"[SHELL]", "PATH"},
	Private: false,
	Action:  actionSimple(watchCommand),
}

func watchCommand(env Env, args []string) (err error) {
	var shellName string

	args = args[1:]
	if len(args) < 1 {
		return fmt.Errorf("a path is required to add to the list of watches")
	}
	if len(args) >= 2 {
		shellName = args[0]
		args = args[1:]
	} else {
		shellName = "bash"
	}

	shell := DetectShell(shellName)

	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", shellName)
	}

	watches := NewFileTimes()
	watchString, ok := env[DIRENV_WATCHES]
	if ok {
		watches.Unmarshal(watchString)
	}

	for idx := range args {
		watches.Update(args[idx])
	}

	e := make(ShellExport)
	e.Add(DIRENV_WATCHES, watches.Marshal())

	os.Stdout.WriteString(shell.Export(e))

	return
}
