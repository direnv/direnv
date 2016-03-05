package main

import (
	"fmt"
	"os"
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

func (self *ExportContext) loadRC() (err error) {
	self.newEnv, err = self.foundRC.Load(self.config, self.oldEnv)

	return
}

func (self *ExportContext) hasRC() bool {
	return self.foundRC != nil || self.loadedRC != nil
}

func (self *ExportContext) updateRC() (err error) {
	self.oldEnv = self.env.Copy()
	var backupDiff *EnvDiff

	if backupDiff, err = self.config.EnvDiff(); err != nil {
		err = fmt.Errorf("EnvDiff() failed: %q", err)
		return
	}
	self.oldEnv = backupDiff.Reverse().Patch(self.env)

	switch {
	case self.foundRC == nil:
		err = self.unloadEnv()
	case self.loadedRC.path != self.foundRC.path:
		err = self.loadRC()
	case self.loadedRC.mtime != self.foundRC.mtime:
		err = self.loadRC()
	}
	return
}

func (self *ExportContext) unloadEnv() (err error) {
	log_status(self.env, "unloading")
	self.newEnv = self.oldEnv.Copy()
	delete(self.newEnv, DIRENV_DIR)
	delete(self.newEnv, DIRENV_MTIME)
	delete(self.newEnv, DIRENV_DIFF)
	return nil
}

func (self *ExportContext) resetEnv() {
	self.newEnv = self.oldEnv.Copy()
	delete(self.oldEnv, DIRENV_DIR)
	delete(self.oldEnv, DIRENV_MTIME)
	delete(self.oldEnv, DIRENV_DIFF)
	if self.foundRC != nil {
		delete(self.newEnv, DIRENV_DIFF)
		self.foundRC.RecordState(self.oldEnv, self.newEnv)
	}
}

func (self *ExportContext) diffString(shell Shell) string {
	oldDiff := self.oldEnv.Diff(self.newEnv)
	if oldDiff.Any() {
		var out []string
		for key, _ := range oldDiff.Prev {
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
	context := ExportContext{env: env}

	var target string

	if len(args) > 1 {
		target = args[1]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("Unknown target shell '%s'", target)
	}

	if err = context.loadConfig(); err != nil {
		return
	}

	context.getRCs()

	if !context.hasRC() {
		return nil
	}

	err = context.updateRC()

	if err != nil {
		context.resetEnv()
	}

	if context.newEnv == nil {
		fmt.Fprintf(os.Stderr, "New env is blank\n\n")
		return nil
	}

	fmt.Print(context.diffString(shell))

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
