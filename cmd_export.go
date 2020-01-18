package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// ExportContext is a sort of state holder struct that is being used to
// record changes before the export finishes.
type ExportContext struct {
	config   *Config
	foundRC  *RC
	loadedRC *RC
	env      Env
	oldEnv   Env
	newEnv   Env
}

func (context *ExportContext) getRCs() {
	context.loadedRC = context.config.LoadedRC()
	context.foundRC = context.config.FindRC()
}

func (context *ExportContext) hasRC() bool {
	return context.foundRC != nil || context.loadedRC != nil
}

func (context *ExportContext) updateRC() (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "update:")

	context.oldEnv = context.env.Copy()
	var backupDiff *EnvDiff

	if backupDiff, err = context.config.EnvDiff(); err != nil {
		err = fmt.Errorf("EnvDiff() failed: %q", err)
		return
	}

	context.oldEnv = backupDiff.Reverse().Patch(context.env)

	logDebug("Determining action:")
	logDebug("foundRC: %#v", context.foundRC)
	logDebug("loadedRC: %#v", context.loadedRC)

	switch {
	case context.foundRC == nil:
		logDebug("no RC found, unloading")
		context.unloadEnv()
	case context.loadedRC == nil:
		logDebug("no RC (implies no DIRENV_DIFF),loading")
		err = context.loadRC()
	case context.loadedRC.path != context.foundRC.path:
		logDebug("new RC, loading")
		err = context.loadRC()
	case context.loadedRC.times.Check() != nil:
		logDebug("file changed, reloading")
		err = context.loadRC()
	default:
		logDebug("no update needed")
	}

	return
}

func (context *ExportContext) loadRC() (err error) {
	context.newEnv, err = context.foundRC.Load(context.config, context.oldEnv)
	return
}

func (context *ExportContext) unloadEnv() {
	logStatus(context.env, "unloading")
	context.newEnv = context.oldEnv.Copy()
	cleanEnv(context.newEnv)
}

func cleanEnv(env Env) {
	env.CleanContext()
}

func (context *ExportContext) diffString(shell Shell) string {
	oldDiff := context.oldEnv.Diff(context.newEnv)
	if oldDiff.Any() {
		var out []string
		for key := range oldDiff.Prev {
			_, ok := oldDiff.Next[key]
			if !ok && !direnvKey(key) {
				out = append(out, "-"+key)
			}
		}

		for key := range oldDiff.Next {
			_, ok := oldDiff.Prev[key]
			if direnvKey(key) {
				continue
			}
			if ok {
				out = append(out, "~"+key)
			} else {
				out = append(out, "+"+key)
			}
		}

		sort.Strings(out)
		if len(out) > 0 {
			logStatus(context.env, "export %s", strings.Join(out, " "))
		}
	}

	diff := context.env.Diff(context.newEnv)
	return diff.ToShell(shell)
}

func exportCommand(env Env, args []string, config *Config) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "export:")
	logDebug("start")
	context := ExportContext{
		env:    env,
		config: config,
	}

	var target string

	if len(args) > 1 {
		target = args[1]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", target)
	}

	logDebug("loading RCs")
	if context.getRCs(); !context.hasRC() {
		return nil
	}

	logDebug("updating RC")
	if err = context.updateRC(); err != nil {
		logDebug("err: %v", err)
	}

	if context.newEnv == nil {
		logDebug("newEnv nil, exiting")
		return
	}

	diffString := context.diffString(shell)
	logDebug("env diff %s", diffString)
	fmt.Print(diffString)

	return
}

// CmdExport is `direnv export $0`
var CmdExport = &Cmd{
	Name:    "export",
	Desc:    "loads an .envrc and prints the diff in terms of exports",
	Args:    []string{"SHELL"},
	Private: true,
	Action:  cmdWithWarnTimeout(actionWithConfig(exportCommand)),
}

func direnvKey(key string) bool {
	return strings.HasPrefix(key, "DIRENV_")
}
