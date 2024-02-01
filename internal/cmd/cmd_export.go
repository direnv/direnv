package cmd

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

func supportedShellFormattedString() string {
	res := "["
	for k := range supportedShellList {
		res += k + ", "
	}
	res = strings.TrimSuffix(res, ", ")
	res += "]"
	return res
}

// CmdExport is `direnv export $0`
var CmdExport = &Cmd{
	Name: "export",
	Desc: `Loads an .envrc or .env and prints the diff in terms of exports.
  Supported SHELL values are: ` + supportedShellFormattedString(),
	Args:    []string{"SHELL"},
	Private: false,
	Action:  cmdWithWarnTimeout(actionWithConfig(exportCommand)),
}

func exportCommand(currentEnv Env, args []string, config *Config) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "export:")
	logDebug("start")

	var target string

	if len(args) > 1 {
		target = args[1]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", target)
	}

	logDebug("loading RCs")
	loadedRC := config.LoadedRC()
	toLoad := findEnvUp(config.WorkDir, config.LoadDotenv)

	if loadedRC == nil && toLoad == "" {
		return
	}

	logDebug("updating RC")
	log.SetPrefix(log.Prefix() + "update:")

	logDebug("Determining action:")
	logDebug("toLoad: %#v", toLoad)
	logDebug("loadedRC: %#v", loadedRC)

	switch {
	case toLoad == "":
		logDebug("no RC found, unloading")
	case loadedRC == nil:
		logDebug("no RC (implies no DIRENV_DIFF),loading")
	case loadedRC.path != toLoad:
		logDebug("new RC, loading")
	case loadedRC.times.Check() != nil:
		logDebug("file changed, reloading")
	default:
		logDebug("no update needed")
		return
	}

	var previousEnv, newEnv Env

	if previousEnv, err = config.Revert(currentEnv); err != nil {
		err = fmt.Errorf("Revert() failed: %w", err)
		logDebug("err: %v", err)
		return
	}

	if toLoad == "" {
		logStatus(currentEnv, "unloading")
		newEnv = previousEnv.Copy()
		newEnv.CleanContext()
	} else {
		newEnv, err = config.EnvFromRC(toLoad, previousEnv)
		if err != nil {
			logDebug("err: %v", err)
			// If loading fails, fall through and deliver a diff anyway,
			// but still exit with an error.  This prevents retrying on
			// every prompt.
		}
		if newEnv == nil {
			// unless of course, the error was in hashing and timestamp loading,
			// in which case we have to abort because we don't know what timestamp
			// to put in the diff!
			return
		}
	}

	if out := diffStatus(previousEnv.Diff(newEnv)); out != "" && !config.HideEnvDiff {
		logStatus(currentEnv, "export %s", out)
	}

	diffString := currentEnv.Diff(newEnv).ToShell(shell)
	logDebug("env diff %s", diffString)
	fmt.Print(diffString)

	return
}

// Return a string of +/-/~ indicators of an environment diff
func diffStatus(oldDiff *EnvDiff) string {
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
		return strings.Join(out, " ")
	}
	return ""
}

func direnvKey(key string) bool {
	return strings.HasPrefix(key, "DIRENV_")
}
