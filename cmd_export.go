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

func (self *ExportContext) loadConfig() (err error) {
	self.config, err = LoadConfig(self.env)
	return
}

func (self *ExportContext) getRCs() {
	self.loadedRC = self.config.LoadedRC()
	self.foundRC = self.config.FindRC()
}

func (self *ExportContext) hasRC() bool {
	return self.foundRC != nil || self.loadedRC != nil
}

func (self *ExportContext) updateRC() (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "update:")

	self.oldEnv = self.env.Copy()
	var backupDiff *EnvDiff

	if backupDiff, err = self.config.EnvDiff(); err != nil {
		err = fmt.Errorf("EnvDiff() failed: %q", err)
		return
	}

	self.oldEnv = backupDiff.Reverse().Patch(self.env)

	log_debug("Determining action:")
	log_debug("foundRC: %#v", self.foundRC)
	log_debug("loadedRC: %#v", self.loadedRC)

	switch {
	case self.foundRC == nil:
		log_debug("no RC found, unloading")
		err = self.unloadEnv()
	case self.loadedRC == nil:
		log_debug("no RC (implies no DIRENV_DIFF),loading")
		err = self.loadRC()
	case self.loadedRC.path != self.foundRC.path:
		log_debug("new RC, loading")
		err = self.loadRC()
	case self.loadedRC.times.Check() != nil:
		log_debug("file changed, reloading")
		err = self.loadRC()
	default:
		log_debug("no update needed")
	}

	return
}

func (self *ExportContext) loadRC() (err error) {
	self.newEnv, err = self.foundRC.Load(self.config, self.oldEnv)
	return
}

func (self *ExportContext) unloadEnv() (err error) {
	log_status(self.env, "unloading")
	self.newEnv = self.oldEnv.Copy()
	cleanEnv(self.newEnv)
	return
}

func (self *ExportContext) resetEnv() {
	self.newEnv = self.oldEnv.Copy()
	cleanEnv(self.oldEnv)
	if self.foundRC != nil {
		delete(self.newEnv, DIRENV_DIFF)
		self.foundRC.RecordState(self.oldEnv, self.newEnv)
	}
}

func cleanEnv(env Env) {
	env.CleanContext()
}

func (self *ExportContext) diffString(shell Shell) string {
	oldDiff := self.oldEnv.Diff(self.newEnv)
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
			log_status(self.env, "export %s", strings.Join(out, " "))
		}
	}

	diff := self.env.Diff(self.newEnv)
	return diff.ToShell(shell)
}

func exportCommand(env Env, args []string) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "export:")
	log_debug("start")
	context := ExportContext{env: env}

	var target string

	if len(args) > 1 {
		target = args[1]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("Unknown target shell '%s'", target)
	}

	log_debug("load config")
	if err = context.loadConfig(); err != nil {
		return
	}

	log_debug("loading RCs")
	if context.getRCs(); !context.hasRC() {
		return nil
	}

	log_debug("updating RC")
	if err = context.updateRC(); err != nil {
		log_debug("err: %v", err)
		context.resetEnv()
	}

	if context.newEnv == nil {
		log_debug("newEnv nil, exiting")
		return nil
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
	Fn:      exportCommand,
}

func direnvKey(key string) bool {
	return strings.HasPrefix(key, "DIRENV_")
}
