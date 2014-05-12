package main

import (
	"fmt"
)

var CmdStatus = &Cmd{
	Name: "status",
	Desc: "prints some debug status information",
	Fn: func(env Env, args []string) error {
		config, err := LoadConfig(env)
		if err != nil {
			return err
		}

		fmt.Println("direnv exec path", config.SelfPath)
		fmt.Println("DIRENV_CONFIG", config.ConfDir)
		fmt.Println("DIRENV_BASH", config.BashPath)

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
	},
}
