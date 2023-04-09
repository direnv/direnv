package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// CmdAllow is `direnv allow [PATH_TO_RC]`
var CmdAllow = &Cmd{
	Name:    "allow",
	Desc:    "Grants direnv permission to load the given .envrc or .env file.",
	Args:    []string{"[PATH_TO_RC]"},
	Aliases: []string{"permit", "grant"},
	Action:  actionWithConfig(cmdAllowAction),
}

var migrationMessage = `
Migrating the allow data to the new location

The allowed .envrc or .env permissions used to be stored in the XDG_CONFIG_HOME. It's
better to keep that folder for user-editable configuration so the data is
being moved to XDG_DATA_HOME.
`

func cmdAllowAction(_ Env, args []string, config *Config) (err error) {
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
			return err
		}
	}

	if _, err = os.Stat(config.AllowDir()); os.IsNotExist(err) {
		oldAllowDir := filepath.Join(config.ConfDir, "allow")
		if _, err = os.Stat(oldAllowDir); err == nil {
			fmt.Println(migrationMessage)

			fmt.Printf("moving %s to %s\n", oldAllowDir, config.AllowDir())
			if err = os.MkdirAll(filepath.Dir(config.AllowDir()), 0755); err != nil {
				return err
			}

			if err = os.Rename(oldAllowDir, config.AllowDir()); err != nil {
				return err
			}

			fmt.Printf("creating a symlink back from %s to %s for back-compat.\n", config.AllowDir(), oldAllowDir)
			if err = os.Symlink(config.AllowDir(), oldAllowDir); err != nil {
				return err
			}
			fmt.Println("")
			fmt.Println("All done, have a nice day!")
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
	return rc.Allow()
}
