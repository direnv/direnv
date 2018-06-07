package main

import (
	"fmt"
	"os"
	"text/template"
)

// HookContext are the variables available during hook template evaluation
type HookContext struct {
	// SelfPath is the unescaped absolute path to direnv
	SelfPath string
}

// `direnv hook $0`
var CmdHook = &Cmd{
	Name: "hook",
	Desc: "Used to setup the shell hook",
	Args: []string{"SHELL"},
	Fn: func(env Env, args []string) (err error) {
		var target string

		if len(args) > 1 {
			target = args[1]
		}

		selfPath, err := os.Executable()
		if err != nil {
			return err
		}

		ctx := HookContext{selfPath}

		shell := DetectShell(target)
		if shell == nil {
			return fmt.Errorf("Unknown target shell '%s'", target)
		}

		hookStr, err := shell.Hook()
		if err != nil {
			return err
		}

		hookTemplate, err := template.New("hook").Parse(hookStr)
		if err != nil {
			return err
		}

		err = hookTemplate.Execute(os.Stdout, ctx)
		if err != nil {
			return err
		}

		return
	},
}
