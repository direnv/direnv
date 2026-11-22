package cmd

import (
	"os"
	"text/template"
)

// HookContext are the variables available during hook template evaluation
type HookContext struct {
	// SelfPath is the unescaped absolute path to direnv
	SelfPath string
}

// BashSelfPath is SelfPath escaped for bash and zsh
func (ctx HookContext) BashSelfPath() string {
	return BashEscape(ctx.SelfPath)
}

// CmdHook is `direnv hook $0`
var CmdHook = &Cmd{
	Name:   "hook",
	Desc:   "Used to setup the shell hook",
	Args:   []string{"SHELL"},
	Action: actionWithConfig(cmdHookAction),
}

func cmdHookAction(_ Env, args []string, config *Config) (err error) {
	var target string

	if len(args) > 1 {
		target = args[1]
	}

	ctx := HookContext{config.SelfPath}

	shell, err := DetectShell(target)
	if err != nil {
		return err
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
