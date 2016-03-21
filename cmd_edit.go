package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		var times *FileTimes
		var foundRC *RC

		defer log.SetPrefix(log.Prefix())
		log.SetPrefix(log.Prefix() + "cmd_edit: ")

		if config, err = LoadConfig(env); err != nil {
			return
		}

		foundRC = config.FindRC()
		if foundRC != nil {
			times = &foundRC.times
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
			log_error("$EDITOR not found.")
			editor = detectEditor(env["PATH"])
			if editor == "" {
				err = fmt.Errorf("Could not find a default editor in the PATH")
				return
			}
		}

		cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("%s %s", editor, rcPath))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return
		}

		foundRC = FindRC(rcPath, config.AllowDir())
		log_debug("foundRC: %#v", foundRC)
		log_debug("times: %#v", times)
		if times != nil {
			log_debug("times.Check(): %#v", times.Check())
		}
		if foundRC != nil && (times == nil || times.Check() != nil) {
			foundRC.Allow()
		}

		return
	},
}

// Utils

var EDITORS = [][]string{
	{"subl", "-w"},
	{"mate", "-w"},
	{"open", "-t", "-W"}, // Opens with the default text editor on mac
	{"nano"},
	{"vim"},
	{"emacs"},
}

func detectEditor(pathenv string) string {
	for _, editor := range EDITORS {
		if _, err := lookPath(editor[0], pathenv); err == nil {
			return strings.Join(editor, " ")
		}
	}
	return ""
}
