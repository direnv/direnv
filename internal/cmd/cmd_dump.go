package cmd

import (
	"fmt"
	"os"
	"strconv"
)

// CmdDump is `direnv dump`
var CmdDump = &Cmd{
	Name:    "dump",
	Desc:    "Used to export the inner bash state at the end of execution",
	Args:    []string{"[SHELL]", "[FILE]"},
	Private: true,
	Action:  actionSimple(cmdDumpAction),
}

func cmdDumpAction(env Env, args []string) (err error) {
	target := "gzenv"
	w := os.Stdout

	if len(args) > 1 {
		target = args[1]
	}

	var filePath string
	if len(args) > 2 {
		filePath = args[2]
	} else {
		filePath = os.Getenv(DIRENV_DUMP_FILE_PATH)
	}

	if filePath != "" {
		if num, err := strconv.Atoi(filePath); err == nil {
			w = os.NewFile(uintptr(num), filePath)
		} else {
			w, err = os.OpenFile(filePath, os.O_WRONLY, 0666)
			if err != nil {
				return err
			}
		}
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", target)
	}

	_, err = fmt.Fprintln(w, shell.Dump(env))

	return
}
