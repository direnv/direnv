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
	Desc:    "Show the data inside of a dump for debugging purposes. Accepts the dump itself or the path to a file containing `direnv dump` output.",
	Args:    []string{"DUMP"},
	Private: true,
	Action:  actionSimple(cmdShowDumpAction),
}

func cmdShowDumpAction(_ Env, args []string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("missing DUMP argument")
	}

	dump, err := dumpArg(args[1])
	if err != nil {
		return err
	}

	var f any
	err = gzenv.Unmarshal(dump, &f)
	if err != nil {
		return err
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	return e.Encode(f)
}

// dumpArg resolves the show_dump argument. When it names an existing file,
// the dump is read from that file (like apply_dump); otherwise the argument
// itself is the dump string.
func dumpArg(arg string) (string, error) {
	if info, err := os.Stat(arg); err == nil && !info.IsDir() {
		data, err := os.ReadFile(arg)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return arg, nil
}
