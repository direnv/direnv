package cmd

import (
	"fmt"
	"os"
)

// CmdWatch is `direnv watch SHELL [PATH...]`
var CmdWatch = &Cmd{
	Name:    "watch",
	Desc:    "Adds a path to the list that direnv watches for changes",
	Args:    []string{"SHELL", "PATH..."},
	Private: true,
	Action:  actionSimple(cmdWatchAction),
}

func cmdWatchAction(env Env, args []string) (err error) {
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
		err = watches.Unmarshal(watchString)
		if err != nil {
			return
		}
	}

	for _, arg := range args[2:] {
		err = watches.Update(arg)
		if err != nil {
			return
		}
	}

	e := make(ShellExport)
	e.Add(DIRENV_WATCHES, watches.Marshal())

	os.Stdout.WriteString(shell.Export(e))

	return
}
