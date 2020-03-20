package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
)

// CmdExport is `direnv export $0`
var CmdExport = &Cmd{
	Name:    "export",
	Desc:    "loads an .envrc and prints the diff in terms of exports",
	Args:    []string{"SHELL"},
	Private: true,
	Action:  cmdWithWarnTimeout(actionWithConfig(actionWithCancel(exportCommand))),
}

func actionWithCancel(fn func(ctx context.Context, env Env, args []string, config *Config) error) actionWithConfig {
	return func(env Env, args []string, config *Config) error {
		ctx, cancel := context.WithCancel(context.Background())

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		// cancel the context on Ctrl-C
		go func() {
			<-c
			cancel()
		}()

		return fn(ctx, env, args, config)
	}
}

func exportCommand(ctx context.Context, env Env, args []string, config *Config) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "export:")
	logDebug("start")

	ec := ExportContext{
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
	if ec.getRCs(); !ec.hasRC() {
		return nil
	}

	logDebug("updating RC")
	if err = ec.updateRC(ctx); err != nil {
		logDebug("err: %v", err)
	}

	if ec.newEnv == nil {
		logDebug("newEnv nil, exiting")
		return
	}

	diffString := ec.diffString(shell)
	logDebug("env diff %s", diffString)
	fmt.Print(diffString)

	return
}

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

func (ec *ExportContext) getRCs() {
	ec.loadedRC = ec.config.LoadedRC()
	ec.foundRC, _ = ec.config.FindRC()
}

func (ec *ExportContext) hasRC() bool {
	return ec.foundRC != nil || ec.loadedRC != nil
}

func (ec *ExportContext) updateRC(ctx context.Context) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "update:")

	ec.oldEnv = ec.env.Copy()

	var backupDiff *EnvDiff
	if backupDiff, err = ec.config.EnvDiff(); err == nil && backupDiff != nil {
		ec.oldEnv = backupDiff.Reverse().Patch(ec.env)
	}
	if err != nil {
		err = fmt.Errorf("EnvDiff() failed: %q", err)
		return
	}

	logDebug("Determining action:")
	logDebug("foundRC: %#v", ec.foundRC)
	logDebug("loadedRC: %#v", ec.loadedRC)

	switch {
	case ec.foundRC == nil:
		logDebug("no RC found, unloading")
		ec.unloadEnv()
	case ec.loadedRC == nil:
		logDebug("no RC (implies no DIRENV_DIFF),loading")
		err = ec.loadRC(ctx)
	case ec.loadedRC.path != ec.foundRC.path:
		logDebug("new RC, loading")
		err = ec.loadRC(ctx)
	case ec.loadedRC.times.Check() != nil:
		logDebug("file changed, reloading")
		err = ec.loadRC(ctx)
	default:
		logDebug("no update needed")
	}

	return
}

func (ec *ExportContext) loadRC(ctx context.Context) (err error) {
	ec.newEnv, err = ec.foundRC.Load(ctx, ec.config, ec.oldEnv)
	return
}

func (ec *ExportContext) unloadEnv() {
	logStatus(ec.env, "unloading")
	ec.newEnv = ec.oldEnv.Copy()
	cleanEnv(ec.newEnv)
}

func cleanEnv(env Env) {
	env.CleanContext()
}

func (ec *ExportContext) diffString(shell Shell) string {
	oldDiff := ec.oldEnv.Diff(ec.newEnv)
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
			logStatus(ec.env, "export %s", strings.Join(out, " "))
		}
	}

	diff := ec.env.Diff(ec.newEnv)
	return diff.ToShell(shell)
}

func direnvKey(key string) bool {
	return strings.HasPrefix(key, "DIRENV_")
}
