package cmd

import (
	"fmt"
	"path/filepath"
)

// CmdStatus is `direnv status`
var CmdStatus = &Cmd{
	Name: "status",
	Desc: "prints some debug status information",
	Action: actionWithConfig(func(env Env, args []string, config *Config) error {
		fmt.Println("direnv exec path", config.SelfPath)
		fmt.Println("DIRENV_CONFIG", config.ConfDir)

		fmt.Println("bash_path", config.BashPath)
		fmt.Println("disable_stdin", config.DisableStdin)
		fmt.Println("warn_timeout", config.WarnTimeout)
		fmt.Println("whitelist.prefix", config.WhitelistPrefix)
		fmt.Println("whitelist.exact", config.WhitelistExact)

		loadedRC := config.LoadedRC()
		foundRC, err := config.FindRC()
		if err != nil {
			return err
		}

		if loadedRC != nil {
			formatRC("Loaded", loadedRC)
		} else {
			fmt.Println("No .envrc or .env loaded")
		}

		if foundRC != nil {
			formatRC("Found", foundRC)
		} else {
			fmt.Println("No .envrc or .env found")
		}

		return nil
	}),
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
