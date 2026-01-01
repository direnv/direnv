package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CmdCheckRequired is `direnv check-required SHELL ENVRC_PATH PATH...`
var CmdCheckRequired = &Cmd{
	Name:    "check-required",
	Desc:    "Checks if required files have been allowed",
	Args:    []string{"SHELL", "ENVRC_PATH", "PATH..."},
	Private: true,
	Action:  actionWithConfig(cmdCheckRequiredAction),
}

func cmdCheckRequiredAction(_ Env, args []string, config *Config) error {
	if len(args) < 2 {
		return fmt.Errorf("a shell name is required")
	}
	if len(args) < 3 {
		return fmt.Errorf("an envrc path is required")
	}
	if len(args) < 4 {
		return fmt.Errorf("at least one file path is required")
	}

	shellName := args[1]
	shell := DetectShell(shellName)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", shellName)
	}

	envrcPath, err := filepath.Abs(args[2])
	if err != nil {
		return fmt.Errorf("failed to resolve envrc path: %w", err)
	}
	envrcDir := filepath.Dir(envrcPath)
	envrcPathHash, err := pathHash(envrcPath)
	if err != nil {
		return fmt.Errorf("failed to hash envrc path: %w", err)
	}

	allowedRequiredDir := filepath.Join(config.AllowedRequiredDir(), envrcPathHash)

	var missingPaths []string

	for _, relPath := range args[3:] {
		// Security validation: must be relative
		if filepath.IsAbs(relPath) {
			fmt.Printf("log_error %s;\n", BashEscape("require_allowed: path must be relative: "+relPath))
			fmt.Println("exit 1;")
			return nil
		}

		// Security validation: no parent traversal
		if strings.Contains(relPath, "..") {
			fmt.Printf("log_error %s;\n", BashEscape("require_allowed: path must not contain '..': "+relPath))
			fmt.Println("exit 1;")
			return nil
		}

		absPath := filepath.Join(envrcDir, relPath)
		hash, err := fileHash(absPath)
		if err != nil {
			// File might not exist or be unreadable
			missingPaths = append(missingPaths, relPath)
			continue
		}

		// Check if the hash exists in the allowed-required directory
		allowedRequiredFile := filepath.Join(allowedRequiredDir, hash)
		if _, err := os.Stat(allowedRequiredFile); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				missingPaths = append(missingPaths, relPath)
			} else {
				return fmt.Errorf("failed to check required file %s: %w", relPath, err)
			}
		}
	}

	if len(missingPaths) > 0 {
		// Export DIRENV_REQUIRED with the missing paths
		e := make(ShellExport)
		e.Add(DIRENV_REQUIRED, strings.Join(missingPaths, ":"))

		exportStr, err := shell.Export(e)
		if err != nil {
			return err
		}
		fmt.Print(exportStr)

		// Output error message
		var fileList string
		if len(missingPaths) == 1 {
			fileList = missingPaths[0] + " requires"
		} else {
			fileList = strings.Join(missingPaths, " and ") + " require"
		}
		fmt.Printf("log_error %s;\n", BashEscape(fileList+" approval. Run 'direnv allow' to approve."))
		fmt.Println("exit 0;")
	}

	return nil
}
