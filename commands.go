package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Cmd struct {
	Name    string
	Desc    string
	Args    []string
	NoWait  bool
	Private bool
	Fn      func(env Env, args []string) error
}

var CmdList []*Cmd

func init() {
	CmdList = []*Cmd{
		CmdAllow,
		CmdDeny,
		CmdDotEnv,
		CmdDump,
		CmdEdit,
		CmdExpandPath,
		CmdExport,
		CmdHelp,
		CmdHook,
		CmdReload,
		CmdStatus,
		CmdStdlib,
		CmdVersion,
	}
}

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
	}

	if command == nil {
		return fmt.Errorf("Command \"%s\" not found", commandPrefix)
	}

	done := make(chan bool, 1)
	if !command.NoWait {
		go func() {
			select {
			case <-done:
				return
			case <-time.After(2 * time.Second):
				fmt.Fprintf(os.Stderr, "direnv(%v) is taking a while to execute. Use CTRL-C to give up.", args)
			}
		}()
	}

	err := command.Fn(env, commandArgs)
	done <- true
	return err
}
