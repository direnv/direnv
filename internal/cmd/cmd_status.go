package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
			fmt.Println("allowed.files", getAllowedFiles(config))
			fmt.Println("denied.files", getDeniedFiles(config))

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

// getAllowedFiles reads all files from the allow directory and returns their original paths
func getAllowedFiles(config *Config) []string {
	var allowed []string
	allowDir := config.AllowDir()
	
	err := filepath.WalkDir(allowDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}
		if d.IsDir() {
			return nil // Skip directories
		}
		
		// Read the file content to get the original path
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}
		
		// The file contains the original path, trim whitespace
		originalPath := strings.TrimSpace(string(content))
		if originalPath != "" {
			allowed = append(allowed, originalPath)
		}
		return nil
	})
	
	if err != nil {
		return []string{} // Return empty slice on error
	}
	
	return allowed
}

// getDeniedFiles reads all files from the deny directory and returns their original paths
func getDeniedFiles(config *Config) []string {
	var denied []string
	denyDir := config.DenyDir()
	
	err := filepath.WalkDir(denyDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors, continue walking
		}
		if d.IsDir() {
			return nil // Skip directories
		}
		
		// Read the file content to get the original path
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}
		
		// The file contains the original path, trim whitespace
		originalPath := strings.TrimSpace(string(content))
		if originalPath != "" {
			denied = append(denied, originalPath)
		}
		return nil
	})
	
	if err != nil {
		return []string{} // Return empty slice on error
	}
	
	return denied
}
