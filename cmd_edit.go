package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// `direnv edit [PATH_TO_RC]`
var CmdEdit = &Cmd{
	Name: "edit",
	Desc: `Opens PATH_TO_RC or the current .envrc into an $EDITOR and allow
  the file to be loaded afterwards.`,
	Args:   []string{"[PATH_TO_RC]"},
	NoWait: true,
	Fn: func(env Env, args []string) (err error) {
		var config *Config
		var rcPath string
		var mtime int64
		var foundRC *RC

		if config, err = LoadConfig(env); err != nil {
			return
		}

		foundRC = config.FindRC()
		if foundRC != nil {
			mtime = foundRC.mtime
		}

		if len(args) > 1 {
			rcPath = args[1]
			fi, _ := os.Stat(rcPath)
			if fi != nil && fi.IsDir() {
				rcPath = filepath.Join(rcPath, ".envrc")
			}
		} else {
			if foundRC == nil {
				return fmt.Errorf(".envrc not found. Use `direnv edit .` to create a new envrc in the current directory.")
			}
			rcPath = foundRC.path
		}

		editor := env["EDITOR"]
		if editor == "" {
			err = fmt.Errorf("$EDITOR not found.")
			return
		}

		cmd := exec.Command(editor, rcPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return
		}

		foundRC = FindRC(rcPath, config.AllowDir())
		if foundRC != nil && foundRC.mtime > mtime {
			foundRC.Allow()
		}

		return
	},
}
