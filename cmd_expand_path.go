package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func expandPath(path, relTo string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(relTo, path))
}

// CmdExpandPath is `direnv expand_path PATH [REL_TO]`
var CmdExpandPath = &Cmd{
	Name:    "expand_path",
	Desc:    "Transforms a PATH to an absolute path to REL_TO or $PWD",
	Args:    []string{"PATH", "[REL_TO]"},
	Private: true,
	Action: actionSimple(func(env Env, args []string) (err error) {
		var path string

		flagset := flag.NewFlagSet(args[0], flag.ExitOnError)
		err = flagset.Parse(args[1:])
		if err != nil {
			return err
		}

		path = flagset.Arg(0)
		if path == "" {
			return fmt.Errorf("PATH missing")
		}

		if !filepath.IsAbs(path) {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			relTo := flagset.Arg(1)
			if relTo == "" {
				relTo = wd
			} else {
				relTo = expandPath(relTo, wd)
			}

			path = expandPath(path, relTo)
		}

		_, err = fmt.Println(path)

		return
	}),
}
