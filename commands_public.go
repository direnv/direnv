package main

import (
	"fmt"
	"os"
)

// `direnv hook $0`
// $0 starts with "-" and go tries to parse it as an argument
//
// This command is public for historical reasons
func Hook(env Env, args []string) (err error) {
	var target string

	if len(args) > 1 {
		target = args[1]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("Unknown target shell '%s'", target)
	}

	fmt.Println(shell.Hook())

	return
}

// `direnv allow [PATH_TO_RC]`
func Allow(env Env, args []string) (err error) {
	var rcPath string
	var config *Config
	if len(args) > 1 {
		rcPath = args[2]
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	if config, err = LoadConfig(env); err != nil {
		return
	}

	rc := FindRC(rcPath, config.AllowDir())
	if rc == nil {
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Allow()
}

// `direnv deny [PATH_TO_RC]`
func Deny(env Env, args []string) (err error) {
	var rcPath string
	var config *Config

	if len(args) > 1 {
		rcPath = args[2]
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	if config, err = LoadConfig(env); err != nil {
		return
	}

	rc := FindRC(rcPath, config.AllowDir())
	if rc == nil {
		return fmt.Errorf(".envrc file not found")
	}
	return rc.Deny()
}

func Status(env Env, args []string) error {
	config, err := LoadConfig(env)
	if err != nil {
		return err
	}

	fmt.Println("DIRENV_LIBEXEC", config.ExecDir)
	fmt.Println("DIRENV_CONFIG", config.ConfDir)

	loadedRC := config.LoadedRC()
	foundRC := config.FindRC()

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
