package main

import (
	"fmt"
	"path/filepath"
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
			formatRC("Loaded", loadedRC)
		} else {
			fmt.Println("No .envrc loaded")
		}

		if foundRC != nil {
			formatRC("Found", foundRC)
		} else {
			fmt.Println("No .envrc found")
		}

		return nil
	},
}

func formatRC(desc string, rc *RC) {
	workDir := filepath.Dir(rc.path)

	fmt.Println(desc, "RC path", rc.path)
	for idx := range *(rc.times.list) {
		fmt.Println(desc, "watch:", (*rc.times.list)[idx].Formatted(workDir))
	}
	fmt.Println(desc, "RC allowed", rc.Allowed())
	fmt.Println(desc, "RC allowPath", rc.allowPath)
}
