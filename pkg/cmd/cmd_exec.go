package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// CmdExec is `direnv exec DIR <COMMAND> ...`
var CmdExec = &Cmd{
	Name:   "exec",
	Desc:   "Executes a command after loading the first .envrc or .env found in DIR",
	Args:   []string{"DIR", "COMMAND", "[...ARGS]"},
	Action: actionWithConfig(cmdExecAction),
}

func cmdExecAction(env Env, args []string, config *Config) (err error) {
	var (
		newEnv      Env
		previousEnv Env
		rcPath      string
		command     string
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

	// Restore pristine environment if needed
	if previousEnv, err = config.Revert(env); err != nil {
		return
	}
	previousEnv.CleanContext()

	// Load the rc
	if toLoad := findEnvUp(rcPath, config.LoadDotenv); toLoad != "" {
		if newEnv, err = config.EnvFromRC(toLoad, previousEnv); err != nil {
			return
		}
	} else {
		newEnv = previousEnv
	}

	var commandPath string
	commandPath, err = lookPath(command, newEnv["PATH"])
	if err != nil {
		err = fmt.Errorf("command '%s' not found on PATH '%s'", command, newEnv["PATH"])
		return
	}

	// #nosec G204
	err = syscall.Exec(commandPath, args, newEnv.ToGoEnv())
	return
}
