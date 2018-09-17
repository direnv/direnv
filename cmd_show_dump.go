package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// `direnv show_dump`
var CmdShowDump = &Cmd{
	Name:    "show_dump",
	Desc:    "Show the data inside of a dump for debugging purposes",
	Args:    []string{"DUMP"},
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		if len(args) < 2 {
			return fmt.Errorf("Missing DUMP argument")
		}

		var f interface{}
		err = unmarshal(args[1], &f)
		if err != nil {
			return err
		}

		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		return e.Encode(f)
	},
}
