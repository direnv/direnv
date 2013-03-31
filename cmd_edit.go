package main

import (
	"fmt"
	"os/exec"
	"os"
)

// `direnv edit [PATH_TO_RC]`
// Opens [PATH_TO_RC] or the current .envrc or a new .envrc into an $EDITOR
// and allow the file to be loaded afterwards.
func Edit(env Env, args []string) (err error) {
	var config *Config
	var rcPath string
	var foundRC *RC

	if config, err = LoadConfig(env); err != nil {
		return
	}

	if len(args) > 1 {
		rcPath = args[1]
	} else {
		foundRC = config.FindRC()
		if foundRC != nil {
			rcPath = foundRC.path
		} else {
			rcPath = ".envrc"
		}
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
	if foundRC != nil {
		foundRC.Allow()
	}

	return
}
