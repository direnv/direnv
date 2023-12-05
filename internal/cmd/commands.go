package cmd

import (
	"fmt"
	"strings"
	"time"
)

type actionSimple func(env Env, args []string) error

func (fn actionSimple) Call(env Env, args []string, _ *Config) error {
	return fn(env, args)
}

type actionWithConfig func(env Env, args []string, config *Config) error

func (fn actionWithConfig) Call(env Env, args []string, config *Config) error {
	var err error
	if config == nil {
		config, err = LoadConfig(env)
		if err != nil {
			return err
		}
	}

	return fn(env, args, config)
}

type action interface {
	Call(env Env, args []string, config *Config) error
}

// Cmd represents a direnv sub-command
type Cmd struct {
	Name    string
	Desc    string
	Args    []string
	Aliases []string
	Private bool
	Action  action
}

// CmdList contains the list of all direnv sub-commands
var CmdList []*Cmd

func init() {
	CmdList = []*Cmd{
		CmdAllow,
		CmdApplyDump,
		CmdShowDump,
		CmdDeny,
		CmdDotEnv,
		CmdDump,
		CmdEdit,
		CmdExec,
		CmdExport,
		CmdFetchURL,
		CmdHelp,
		CmdHook,
		CmdPrune,
		CmdReload,
		CmdStatus,
		CmdStdlib,
		CmdVersion,
		CmdWatch,
		CmdWatchDir,
		CmdWatchList,
		CmdWatchPrint,
		CmdCurrent,
	}
}

func cmdWithWarnTimeout(fn action) action {
	return actionWithConfig(func(env Env, args []string, config *Config) (err error) {
		// Disable warning if WarnTimeout is <= 0
		if config.WarnTimeout <= 0 {
			return fn.Call(env, args, config)
		}

		done := make(chan bool, 1)
		go func() {
			select {
			case <-done:
				return
			case <-time.After(config.WarnTimeout):
				logError("(%v) is taking a while to execute. Use CTRL-C to give up.", args)
			}
		}()

		err = fn.Call(env, args, config)
		done <- true
		return err
	})
}

// CommandsDispatch is called by the main() function to dispatch to a sub-command
func CommandsDispatch(env Env, args []string) error {
	var command *Cmd
	var commandName string
	var commandPrefix string
	var commandArgs []string

	if len(args) < 2 {
		commandName = "help"
		commandPrefix = args[0]
		commandArgs = []string{}
	} else {
		commandName = args[1]
		commandPrefix = strings.Join(args[0:2], " ")
		commandArgs = append([]string{commandPrefix}, args[2:]...)
	}

	for _, cmd := range CmdList {
		if cmd.Name == commandName {
			command = cmd
			break
		}
		for _, alias := range cmd.Aliases {
			if alias == commandName {
				command = cmd
			}
		}
	}

	if command == nil {
		return fmt.Errorf("command \"%s\" not found", commandPrefix)
	}

	return command.Action.Call(env, commandArgs, nil)
}
