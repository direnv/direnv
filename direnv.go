package main

import (
	"fmt"
	"os"
)

var privateCommands = CommandDispatcher(map[string]Command{
	"dump":    Dump,
	"diff":    Diff,
	"export":  Export,
	"default": Usage,
})

var stdlibCommands = CommandDispatcher(map[string]Command{
	"default":     Stdlib,
	"expand_path": ExpandPath,
	"mtime":       FileMtime,
	"hash":        FileHash,
})

var publicCommands = CommandDispatcher(map[string]Command{
	"status":  TODO,
	"allow":   TODO,
	"deny":    TODO,
	"switch":  TODO,
	"hook":    Hook,
	"private": privateCommands,
	"stdlib":  stdlibCommands,
	"help":    Usage,
	"default": Usage,
})

func TODO(args []string) error {
	fmt.Println("TODO")
	return nil
}

func main() {
	env := GetEnv()
	context := LoadContext(env)
	CONTEXT = context

	if false {
		fmt.Println(context)
	}

	if err := publicCommands(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
