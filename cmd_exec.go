package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// `direnv exec [DIR] <COMMAND> ...`
var CmdExec = &Cmd{
	Name: "exec",
	Desc: "Executes a command after loading the first .envrc found in DIR",
	Args: []string{"[DIR]", "COMMAND", "[...ARGS]"},
	Fn: func(env Env, args []string) (err error) {
		var (
			backupDiff *EnvDiff
			config     *Config
			newEnv     Env
			rcPath     string
			command    string
		)

		if len(args) < 2 {
			return fmt.Errorf("missing DIR and COMMAND arguments")
		}

		rcPath = filepath.Clean(args[1])
		fi, err := os.Stat(rcPath)
		if err != nil {
			return
		}

		if fi.IsDir() {
			if len(args) < 3 {
				return fmt.Errorf("missing COMMAND argument")
			}
			command = args[2]
			args = args[2:]
		} else {
			command = rcPath
			rcPath = filepath.Dir(rcPath)
			args = args[1:]
		}

		if config, err = LoadConfig(env); err != nil {
			return
		}

		rc := FindRC(rcPath, config.AllowDir())
		if rc == nil {
			return fmt.Errorf(".envrc not found")
		}

		// Restore pristine environment if needed
		if backupDiff, err = config.EnvDiff(); err == nil {
			backupDiff.Reverse().Patch(env)
		}
		delete(env, DIRENV_DIR)
		delete(env, DIRENV_MTIME)
		delete(env, DIRENV_DIFF)

		// Load the rc
		if newEnv, err = rc.Load(config, env); err != nil {
			return
		}

		command, err = lookPath(command, newEnv["PATH"])
		if err != nil {
			return
		}

		err = syscall.Exec(command, args, newEnv.ToGoEnv())
		return
	},
}

// Similar to os/exec.LookPath except we pass in the PATH
func lookPath(file string, pathenv string) (string, error) {
	if strings.Contains(file, "/") {
		err := findExecutable(file)
		if err == nil {
			return file, nil
		}
		return "", err
	}
	if pathenv == "" {
		return "", errNotFound
	}
	for _, dir := range strings.Split(pathenv, ":") {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := dir + "/" + file
		if err := findExecutable(path); err == nil {
			return path, nil
		}
	}
	return "", errNotFound
}

// ErrNotFound is the error resulting if a path search failed to find an executable file.
var errNotFound = errors.New("executable file not found in $PATH")

func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}
