package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// CompleteContext has variables for the "complete" template evaluation
type CompleteContext struct {
	// SelfPath is the unescaped absolute path to direnv
	SelfPath string
	// CmdList is the list of commands to complete
	CmdList []*Cmd
	// CmdMap is a name->command lookup
	CmdMap map[string]*Cmd
}

// CmdComplete is `direnv complete $0`
var CmdComplete = &Cmd{
	Name:   "complete",
	Desc:   "Generate completion for given shell",
	Args:   []string{"SHELL"},
	Action: actionSimple(cmdCompleteAction),
}

// TODO: move this to Shell.Complete() to add support for other shells
const completeStr = `# To add autocomplete for direnv run:
# {{.SelfPath}} complete fish >~/.config/fish/completions/direnv.fish

complete -c direnv -l help -d "{{.CmdMap.help.Desc}}"
complete -c direnv -l version -d "{{.CmdMap.version.Desc}}"

# Subcommands
{{- range $cmd := .CmdList}}
complete -c direnv -x -n "__fish_use_subcommand" -a {{$cmd.Name}} -d "{{$cmd.Desc}}"
{{- end}}

# Arguments
complete -c direnv -x -n "__fish_seen_subcommand_from {{.CmdMap.hook.Name}} {{.CmdMap.complete.Name}}" -a "bash elvish fish tcsh zsh"
`

func cmdCompleteAction(env Env, args []string) (err error) {
	var target string

	if len(args) > 1 {
		target = args[1]
	}
	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", target)
	}
	if shell != Fish {
		return fmt.Errorf("`complete` not supported for `%s`", target)
	}

	selfPath, err := os.Executable()
	if err != nil {
		return err
	}
	// Convert Windows path to Unix if needed
	selfPath = strings.Replace(selfPath, "\\", "/", -1)

	ctx := CompleteContext{
		SelfPath: selfPath,
		CmdList: CmdList,
		CmdMap: make(map[string]*Cmd),
	}
	// Fill in the name-lookup
	for _, cmd := range CmdList {
		ctx.CmdMap[cmd.Name] = cmd
	}

	completeTemplate, err := template.New("complete").Parse(completeStr)
	if err != nil {
		return err
	}

	err = completeTemplate.Execute(os.Stdout, ctx)
	if err != nil {
		return err
	}

	return
}
