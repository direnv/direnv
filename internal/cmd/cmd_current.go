package cmd

import (
	"errors"
)

// CmdCurrent is `direnv current`
var CmdCurrent = &Cmd{
	Name:    "current",
	Desc:    "Reports whether direnv's view of a file is current (or stale)",
	Args:    []string{"PATH"},
	Private: true,
	Action:  actionSimple(cmdCurrentAction),
}

func cmdCurrentAction(env Env, args []string) (err error) {
	if len(args) < 2 {
		err = errors.New("missing PATH argument")
		return
	}

	path := args[1]
	watches := NewFileTimes()
	watchString, ok := env[DIRENV_WATCHES]
	if ok {
		err = watches.Unmarshal(watchString)
		if err != nil {
			return
		}
	}

	err = watches.CheckOne(path)

	return
}
