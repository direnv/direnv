package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CmdEdit is `direnv edit [PATH_TO_RC]`
var CmdEdit = &Cmd{
	Name: "edit",
	Desc: `Opens PATH_TO_RC or the current .envrc or .env into an $EDITOR and allow
  the file to be loaded afterwards.`,
	Args:   []string{"[PATH_TO_RC]"},
	Action: actionWithConfig(cmdEditAction),
}

func cmdEditAction(env Env, args []string, config *Config) (err error) {
	var rcPath string
	var times *FileTimes
	var foundRC *RC

	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "cmd_edit: ")

	foundRC, err = config.FindRC()
	if err != nil {
		return err
	}
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
			return fmt.Errorf(".envrc or .env not found. Use `direnv edit .` to create a new .envrc in the current directory")
		}
		rcPath = foundRC.path
	}

	editor := env["EDITOR"]
	if editor == "" {
		logError("$EDITOR not found.")
		editor = detectEditor(env["PATH"])
		if editor == "" {
			err = fmt.Errorf("could not find a default editor in the PATH")
			return
		}
	}

	run := fmt.Sprintf("%s %s", editor, BashEscape(rcPath))

	// G204: Subprocess launched with function call as argument or cmd arguments
	// #nosec
	cmd := exec.Command(config.BashPath, "-c", run)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return
	}

	foundRC, err = FindRC(rcPath, config)
	logDebug("foundRC: %#v", foundRC)
	logDebug("times: %#v", times)
	if times != nil {
		logDebug("times.Check(): %#v", times.Check())
	}
	if err == nil && foundRC != nil && (times == nil || times.Check() != nil) {
		err = foundRC.Allow()
	}

	return
}

// Utils

// Editors contains a list of known editors and how to start them.
var Editors = [][]string{
	{"editor"},
	{"subl", "-w"},
	{"mate", "-w"},
	{"open", "-t", "-W"}, // Opens with the default text editor on mac
	{"nano"},
	{"vim"},
	{"emacs"},
}

func detectEditor(pathenv string) string {
	for _, editor := range Editors {
		if _, err := lookPath(editor[0], pathenv); err == nil {
			return strings.Join(editor, " ")
		}
	}
	return ""
}
