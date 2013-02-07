package main

import (
	"fmt"
	"strings"
)

type Command func(env Env, args []string) error

func CommandDispatcher(commands map[string]Command) Command {
	return func(env Env, args []string) error {
		var command Command
		var commandName string
		var commandPrefix string
		var commandArgs []string

		if len(args) < 2 {
			commandName = "default"
			commandPrefix = args[0]
			commandArgs = []string{}
		} else {
			commandName = args[1]
			commandPrefix = strings.Join(args[0:2], " ")
			commandArgs = append([]string{commandPrefix}, args[2:]...)
		}

		command = commands[commandName]

		if command == nil {
			command = commandNotFound(commandPrefix)
		}

		return command(env, commandArgs)
	}
}

func commandNotFound(commandPrefix string) Command {
	return func(env Env, args []string) error {
		return fmt.Errorf("Command \"%s\" not found", commandPrefix)
	}
}
