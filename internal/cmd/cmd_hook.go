package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// HookContext are the variables available during hook template evaluation
type HookContext struct {
	// SelfPath is the unescaped absolute path to direnv
	SelfPath string
}

// CmdHook is `direnv hook $0`
var CmdHook = &Cmd{
	Name:   "hook",
	Desc:   "Used to setup the shell hook",
	Args:   []string{"SHELL"},
	Action: actionSimple(cmdHookAction),
}

func cmdHookAction(_ Env, args []string) (err error) {
	var target string

	if len(args) > 1 {
		target = args[1]
	}

	selfPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Convert Windows path if needed
	selfPath = strings.Replace(selfPath, "\\", "/", -1)
	ctx := HookContext{selfPath}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", target)
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
}
