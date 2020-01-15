package main

import (
	"fmt"
	"os"
)

var CmdWatch = &Cmd{
	Name:    "watch",
	Desc:    "Adds a path to the list that direnv watches for changes",
	Args:    []string{"SHELL", "PATH"},
	Private: true,
	Action:  actionSimple(watchCommand),
}

func watchCommand(env Env, args []string) (err error) {
	var shellName string

	if len(args) < 2 {
		return fmt.Errorf("a path is required to add to the list of watches")
	}
	if len(args) >= 2 {
		shellName = args[1]
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

	for _, arg := range args[2:] {
		watches.Update(arg)
	}

	e := make(ShellExport)
	e.Add(DIRENV_WATCHES, watches.Marshal())

	os.Stdout.WriteString(shell.Export(e))

	return
}
