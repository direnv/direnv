package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// CmdExec is `direnv exec DIR <COMMAND> ...`
var CmdExec = &Cmd{
	Name:   "exec",
	Desc:   "Executes a command after loading the first .envrc found in DIR",
	Args:   []string{"DIR", "COMMAND", "[...ARGS]"},
	Action: actionWithConfig(cmdExecAction),
}

func cmdExecAction(env Env, args []string, config *Config) (err error) {
	var (
		backupDiff *EnvDiff
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

	rc, err := FindRC(rcPath, config)
	if err != nil {
		return
	}

	// Set the default physical OS path needed for command invocation, as it
	// may be different than the logical PATH manipulated by the stdlib
	env[DIRENV_PLATFORM_PATH] = env["PATH"]

	// Restore pristine environment if needed
	if backupDiff, err = config.EnvDiff(); err == nil && backupDiff != nil {
		env = backupDiff.Reverse().Patch(env)
	}
	env.CleanContext()

	// Load the rc
	if rc != nil {
		if newEnv, err = rc.Load(context.Background(), config, env); err != nil {
			return
		}
	} else {
		newEnv = env
	}

	var commandPath string
	commandPath, err = lookPath(command, newEnv[DIRENV_PLATFORM_PATH])
	if err != nil {
		err = fmt.Errorf("command '%s' not found on PATH '%s'", command, newEnv[DIRENV_PLATFORM_PATH])
		return
	}

	err = syscall.Exec(commandPath, args, newEnv.ToGoEnv())
	return
}
