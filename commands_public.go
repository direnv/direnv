package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// `direnv hook $0`
// $0 starts with "-" and go tries to parse it as an argument
//
// This command is public for historical reasons
func Hook(env Env, args []string) (err error) {
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

	// $0 starts with "-" but Base doesn't care
	target = filepath.Base(target)

	switch target {
	case "bash":
		fmt.Print(HOOK_BASH)
	case "zsh":
		fmt.Print(HOOK_ZSH)
	default:
		return fmt.Errorf("Unknown target shell '%s'", target)
	}

	return
}

// `direnv allow [PATH_TO_RC]`
func Allow(env Env, args []string) (err error) {
	var rcPath string
	var context *Context
	if len(args) > 1 {
		rcPath = args[2]
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	if context, err = LoadContext(env); err != nil {
		return
	}

	rc := FindRC(rcPath, context.AllowDir())
	if rc == nil {
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Allow()
}

// `direnv deny [PATH_TO_RC]`
func Deny(env Env, args []string) (err error) {
	var rcPath string
	var context *Context

	if len(args) > 1 {
		rcPath = args[2]
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	if context, err = LoadContext(env); err != nil {
		return
	}

	rc := FindRC(rcPath, context.AllowDir())
	if rc == nil {
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Deny()
}

func Switch(env Env, args []string) error {
	return fmt.Errorf("Woops ! This should be handled by the shell wrapper")
}

func Status(env Env, args []string) error {
	context, err := LoadContext(env)
	if err != nil {
		return err
	}

	fmt.Println("DIRENV_LIBEXEC", context.ExecDir)
	fmt.Println("DIRENV_CONFIG", context.ConfDir)

	loadedRC := context.LoadedRC()
	foundRC := context.FoundRC()

	if loadedRC != nil {
		fmt.Println("Loaded RC path", loadedRC.path)
		fmt.Println("Loaded RC mtime", loadedRC.mtime)
		fmt.Println("Loaded RC allowed", loadedRC.Allowed())
		fmt.Println("Loaded RC allowPath", loadedRC.allowPath)
	} else {
		fmt.Println("No .envrc loaded")
	}

	if foundRC != nil {
		fmt.Println("Found RC path", foundRC.path)
		fmt.Println("Found RC mtime", foundRC.mtime)
		fmt.Println("Found RC allowed", foundRC.Allowed())
		fmt.Println("Found RC allowPath", foundRC.allowPath)
	} else {
		fmt.Println("No .envrc found")
	}

	return nil
}
