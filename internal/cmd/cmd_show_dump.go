package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/direnv/direnv/v2/gzenv"
)

// CmdShowDump is `direnv show_dump`
var CmdShowDump = &Cmd{
	Name:    "show_dump",
	Desc:    "Show the data inside of a dump for debugging purposes",
	Args:    []string{"DUMP"},
	Private: true,
	Action:  actionSimple(cmdShowDumpAction),
}

func cmdShowDumpAction(_ Env, args []string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("missing DUMP argument")
	}

	var f interface{}
	err = gzenv.Unmarshal(args[1], &f)
	if err != nil {
		return err
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	return e.Encode(f)
}
