package cmd

import (
	"fmt"
	"path/filepath"
)

// CmdExpandPath is `direnv expand_path <path>`
var CmdExpandPath = &Cmd{
	Name:    "expand_path",
	Desc:    "Transforms path to an absolute path",
	Args:    []string{"<rel_path>", "[<relative_to>]"},
	Private: true,
	Action:  actionWithConfig(cmdExpandPathAction),
}

func cmdExpandPathAction(env Env, args []string, config *Config) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("expecting at least one argument, got %d", len(args)-1)
	}
	var path string
	if path, err = filepath.Abs(args[1]); err != nil {
		return err
	}
	// Build the path relative to
	if len(args) > 2 {
		var base string
		if base, err = filepath.Abs(args[2]); err != nil {
			return err
		}
		if path, err = filepath.Rel(base, path); err != nil {
			return err
		}
	}

	fmt.Println(path)
	return nil
}
