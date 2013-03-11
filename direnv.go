package main

import (
	"fmt"
	"os"
)

var privateCommands = CommandDispatcher(map[string]Command{
	"dump":        Dump,
	"expand_path": ExpandPath,
	"export":      Export,
	"load":        Load,
})

var publicCommands = CommandDispatcher(map[string]Command{
	"allow":   Allow,
	"default": TODO,
	"deny":    Deny,
	"help":    TODO,
	"hook":    Hook,
	"private": privateCommands,
	"status":  Status,
	"switch":  Switch,
})

func TODO(env Env, args []string) error {
	fmt.Println("TODO")
	return nil
}

func main() {
	env := GetEnv()

	if err := publicCommands(env, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
