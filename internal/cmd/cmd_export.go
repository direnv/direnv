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

const (
	exportContextNoProcessMarker = "no_process_marker"
	exportContextExit            = "exit"
)

// CmdExport is `direnv export $0`
var CmdExport = &Cmd{
	Name: "export",
	Desc: `Loads an .envrc or .env and prints the diff in terms of exports.
  Supported SHELL values are: ` + supportedShellFormattedString(),
	Args:    []string{"SHELL", "[--context]", "[" + exportContextNoProcessMarker + " | " + exportContextExit + "]"},
	Private: false,
	Action:  cmdWithWarnTimeout(actionWithConfig(exportCommand)),
}

func exportCommand(currentEnv Env, args []string, config *Config) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "export:")
	logDebug("start")

	var target, exportContext string

	if len(args) > 1 {
		target = args[1]
	}
	if len(args) > 3 {
		exportContext = args[3]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", target)
	}
	hookableShell, shellIsHookable := shell.(HookableShell)

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
	case currentEnv[DIRENV_REQUIRED] != "":
		// Force reload if required files were pending approval.
		// The approval status might have changed even if file times haven't.
		logDebug("required files pending, reloading")
	case shellIsHookable && exportContext == exportContextExit:
		logDebug("the shell is exiting, running hooks")
	case shellIsHookable && exportContext == exportContextNoProcessMarker:
		logDebug("a new shell was started, running hooks")
	default:
		logDebug("no update needed")
		return
	}

	var currentEnvState *State
	currentEnvState, err = UnmarshalState(currentEnv[DIRENV_STATE])
	if err != nil {
		return err
	}

	if shellIsHookable && loadedRC != nil && exportContext != "" {
		hooks := map[string]string{}
		var setProcessMarker *bool
		switch exportContext {
		// t_exit
		case exportContextExit:
			addHookIfPresent(currentEnvState, hookableShell, hooks, HOOK_UNLOAD)
		// t_exec, t_subshell
		case exportContextNoProcessMarker:
			setProcessMarker = getBoolPointer(true)
			addHookIfPresent(currentEnvState, hookableShell, hooks, HOOK_LOAD)
		}

		return printExportWithHooks(hookableShell, nil, hooks, setProcessMarker)
	}

	var previousEnv, newEnv Env

	if previousEnv, err = config.Revert(currentEnv); err != nil {
		err = fmt.Errorf("Revert() failed: %w", err)
		logDebug("err: %v", err)
		return
	}

	var loadNewEnvErr error
	var newEnvState *State
	if toLoad == "" {
		logStatus(config, "unloading")
		newEnv = previousEnv.Copy()
		newEnv.CleanContext()
	} else {
		newEnv, loadNewEnvErr = config.EnvFromRC(toLoad, previousEnv)
		if loadNewEnvErr != nil {
			logDebug("err: %v", loadNewEnvErr)
			// If loading fails, fall through and deliver a diff anyway,
			// but still exit with an error.  This prevents retrying on
			// every prompt.
		}
		if newEnv == nil {
			// unless of course, the error was in hashing and timestamp loading,
			// in which case we have to abort because we don't know what timestamp
			// to put in the diff!
			return loadNewEnvErr
		}

		newEnvState = MakeState(newEnv)
		newEnv[DIRENV_STATE] = newEnvState.Marshal()
	}

	if out := diffStatus(previousEnv.Diff(newEnv)); out != "" && !config.HideEnvDiff {
		logStatus(config, "export %s", out)
	}

	if shellIsHookable {
		shellExport := currentEnv.Diff(newEnv).ToShellExport()
		hooks := map[string]string{}
		var setProcessMarker *bool
		if toLoad == "" {
			// t_cd_outside_direnv
			setProcessMarker = getBoolPointer(false)
			addHookIfPresent(currentEnvState, hookableShell, hooks, HOOK_UNLOAD)
		} else {
			if loadedRC == nil {
				// t_cd_to_direnv
				setProcessMarker = getBoolPointer(true)
				addHookIfPresent(newEnvState, hookableShell, hooks, HOOK_LOAD)
			} else {
				// t_change_watched_file, t_cd_to_a_different_direnv, t_block, t_allow
				addHookIfPresent(currentEnvState, hookableShell, hooks, HOOK_UNLOAD)
				addHookIfPresent(newEnvState, hookableShell, hooks, HOOK_LOAD)
			}
		}

		err = printExportWithHooks(hookableShell, shellExport, hooks, setProcessMarker)
		if err != nil {
			return err
		}
	} else {
		err = printExport(shell, currentEnv.Diff(newEnv))
		if err != nil {
			return err
		}
	}

	return loadNewEnvErr
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

func printExport(shell Shell, envDiff *EnvDiff) error {
	diffString, diffErr := envDiff.ToShell(shell)
	if diffErr != nil {
		return fmt.Errorf("ToShell() failed: %w", diffErr)
	}

	logDebug("env diff %s", diffString)

	fmt.Print(diffString)

	return nil
}

func printExportWithHooks(hookableShell HookableShell, shellExport ShellExport, hooks map[string]string, setProcessMarker *bool) error {
	diffString, diffErr := hookableShell.ExportWithHooks(shellExport, hooks, setProcessMarker)
	if diffErr != nil {
		return fmt.Errorf("ExportWithHooks() failed: %w", diffErr)
	}

	logDebug("env diff %s", diffString)

	fmt.Print(diffString)

	return nil
}

func addHookIfPresent(state *State, hookableShell HookableShell, hooks map[string]string, hookName string) {
	hook := state.Hooks.Get(hookName, hookableShell.Name())
	if hook != "" {
		hooks[hookName] = hook
	}
}

// TODO: When the minimum Go version for the project reaches 1.26, we can instead use `new(<bool>)`
func getBoolPointer(b bool) *bool {
	return &b
}
