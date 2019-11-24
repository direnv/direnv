package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

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

	log_debug("Determining action:")
	log_debug("foundRC: %#v", context.foundRC)
	log_debug("loadedRC: %#v", context.loadedRC)

	switch {
	case context.foundRC == nil:
		log_debug("no RC found, unloading")
		err = context.unloadEnv()
	case context.loadedRC == nil:
		log_debug("no RC (implies no DIRENV_DIFF),loading")
		err = context.loadRC()
	case context.loadedRC.path != context.foundRC.path:
		log_debug("new RC, loading")
		err = context.loadRC()
	case context.loadedRC.times.Check() != nil:
		log_debug("file changed, reloading")
		err = context.loadRC()
	default:
		log_debug("no update needed")
	}

	return
}

func (context *ExportContext) loadRC() (err error) {
	context.newEnv, err = context.foundRC.Load(context.config, context.oldEnv)
	return
}

func (context *ExportContext) unloadEnv() (err error) {
	log_status(context.env, "unloading")
	context.newEnv = context.oldEnv.Copy()
	cleanEnv(context.newEnv)
	return
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
			log_status(context.env, "export %s", strings.Join(out, " "))
		}
	}

	diff := context.env.Diff(context.newEnv)
	return diff.ToShell(shell)
}

func exportCommand(env Env, args []string, config *Config) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "export:")
	log_debug("start")
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

	log_debug("loading RCs")
	if context.getRCs(); !context.hasRC() {
		return nil
	}

	log_debug("updating RC")
	if err = context.updateRC(); err != nil {
		log_debug("err: %v", err)
	}

	if context.newEnv == nil {
		log_debug("newEnv nil, exiting")
		return
	}

	diffString := context.diffString(shell)
	log_debug("env diff %s", diffString)
	fmt.Print(diffString)

	return
}

// `direnv export $0`
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
