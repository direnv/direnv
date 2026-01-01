package cmd

import (
	"fmt"
)

// CmdWatchPrint is `direnv watch-print`
var CmdWatchPrint = &Cmd{
	Name:    "watch-print",
	Desc:    "prints the watched paths",
	Args:    []string{"[--null]"},
	Private: true,
	Action:  actionSimple(cmdWatchPrintAction),
}

func cmdWatchPrintAction(env Env, args []string) (err error) {
	watches := NewFileTimes()
	watchString, ok := env[DIRENV_WATCHES]
	separator := '\n'
	if len(args) > 1 && args[1] == "--null" {
		separator = 0
	}

	if ok {
		err = watches.Unmarshal(watchString)
		if err != nil {
			return
		}
	}

	for _, watch := range *watches.list {
		fmt.Printf("%s%c", watch.Path, separator)
	}

	return
}
