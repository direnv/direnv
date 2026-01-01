package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

// CmdStatus is `direnv status`
var CmdStatus = &Cmd{
	Name: "status",
	Desc: "Prints some debug status information",
	Args: []string{"[--json]"},
	Action: actionWithConfig(func(_ Env, args []string, config *Config) error {
		if len(args) > 1 && (args[1] == "-json" || args[1] == "--json") {
			loadedRC := config.LoadedRC()
			foundRC, err := config.FindRC()
			if err != nil {
				return err
			}
			jsonOutput := map[string]interface{}{
				"config": map[string]string{
					"SelfPath":  config.SelfPath,
					"ConfigDir": config.ConfDir,
				},
				"state": map[string]interface{}{},
			}
			if loadedRC != nil {
				jsonOutput["state"].(map[string]interface{})["loadedRC"] = map[string]interface{}{
					"path":    loadedRC.path,
					"allowed": loadedRC.Allowed(),
				}
			} else {
				jsonOutput["state"].(map[string]interface{})["loadedRC"] = nil
			}
			if foundRC != nil {
				jsonOutput["state"].(map[string]interface{})["foundRC"] = map[string]interface{}{
					"path":    foundRC.path,
					"allowed": foundRC.Allowed(),
				}
			} else {
				jsonOutput["state"].(map[string]interface{})["foundRC"] = nil
			}
			jsonBytes, err := json.MarshalIndent(jsonOutput, "", "  ")
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Println(string(jsonBytes))
		} else {
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
