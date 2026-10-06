package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// CmdDeny is `direnv deny [PATH_TO_RC]`
var CmdDeny = &Cmd{
	Name:    "block",
	Desc:    "Revokes the authorization of a given .envrc or .env file.",
	Args:    []string{"[PATH_TO_RC]"},
	Aliases: []string{"deny", "disallow", "revoke"},
	Action:  actionWithConfig(cmdDenyAction),
}

func cmdDenyAction(_ Env, args []string, config *Config) (err error) {
	var rcPath string

	if len(args) > 1 {
		if rcPath, err = filepath.Abs(args[1]); err != nil {
			return err
		}
		if rcPath, err = filepath.EvalSymlinks(rcPath); err != nil {
			return err
		}
	} else {
		if rcPath, err = os.Getwd(); err != nil {
			return
		}
	}

	rc, err := FindRC(rcPath, config)
	if err != nil {
		return err
	} else if rc == nil {
		if config.LoadDotenv {
			return fmt.Errorf(".envrc or .env file not found")
		}
		return fmt.Errorf(".envrc file not found")
	}

	// Remove required files for this .envrc
	if err = removeAllowedRequiredFiles(rc.Path(), config); err != nil {
		return err
	}

	return rc.Deny()
}

func removeAllowedRequiredFiles(rcPath string, config *Config) error {
	envrcPathHash, err := pathHash(rcPath)
	if err != nil {
		return fmt.Errorf("failed to hash envrc path: %w", err)
	}

	allowedRequiredDir := filepath.Join(config.AllowedRequiredDir(), envrcPathHash)

	// Remove the entire allowed-required directory for this .envrc
	if err := os.RemoveAll(allowedRequiredDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove allowed-required files: %w", err)
	}

	return nil
}
