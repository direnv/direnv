package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NOTE: direnv hook $0
// $0 starts with "-" and go tries to parse it as an argument
//
// This command is public for historical reasons
func Hook(args []string) (err error) {
	var target string

	if len(args) > 1 {
		target = args[1]
	} else {
		// Try to find out the shell on Linux systems
		ppid := os.Getppid()
		data, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", ppid))
		if err != nil {
			return fmt.Errorf("Please specify a target shell")
		}

		target = string(data)
	}

	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	target = filepath.Base(target)

	switch target {
	case "bash":
		fmt.Println("PROMPT_COMMAND=\"eval \\`direnv export\\`;$PROMPT_COMMAND")
	case "zsh":
		fmt.Println("direnv_hook() { eval `direnv export` }; [[ -z $precmd_functions ]] && precmd_functions=(); precmd_functions=($precmd_functions direnv_hook)")
	default:
		return fmt.Errorf("Unknown target shell '%s'", target)
	}

	return
}

func Usage(args []string) error {
	fmt.Println("HI !")
	return nil
}
