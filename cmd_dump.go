package main

import (
	"fmt"
	"io"
	"os"
)

// `direnv dump`
var CmdDump = &Cmd{
	Name:    "dump",
	Desc:    "Used to export the inner bash state at the end of execution",
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		output := os.NewFile(uintptr(3), "DIRENV_DUMP_FD")

		_, err = io.WriteString(output, env.Filtered().Serialize())
		if err != nil {
			return fmt.Errorf("dump: %v", err)
		}

		_, err = output.Write([]byte{0, '\n'})
		if err != nil {
			return fmt.Errorf("dump: %v", err)
		}

		return
	},
}
