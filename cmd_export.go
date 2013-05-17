package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// `direnv export $0`
var CmdExport = &Cmd{
	Name:    "export",
	Desc:    "loads an .envrc and prints the diff in terms of exports",
	Args:    []string{"SHELL"},
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		var oldEnv Env = env.Filtered()
		var newEnv Env
		var loadedRC *RC
		var foundRC *RC
		var config *Config
		var target string

		if len(args) > 1 {
			target = args[1]
		}

		shell := DetectShell(target)
		if shell == nil {
			return fmt.Errorf("Unknown target shell '%s'", target)
		}

		if config, err = LoadConfig(env); err != nil {
			return
		}

		loadedRC = config.LoadedRC()
		foundRC = config.FindRC()

		//fmt.Fprintf(os.Stderr, "%v %v\n", loadedRC, foundRC)

		if loadedRC == nil {
			if foundRC == nil {
				// We're done here.
				return nil
			}

			newEnv, err = loadRC(foundRC, config, oldEnv)
		} else {
			var backupEnv Env
			if backupEnv, err = config.EnvBackup(); err != nil {
				err = fmt.Errorf("EnvBackup() failed: %q", err)
				goto error
			}
			oldEnv = backupEnv.Filtered()
			if foundRC == nil {
				fmt.Fprintf(os.Stderr, "direnv: unloading\n")
				newEnv = oldEnv
			} else if loadedRC.path != foundRC.path {
				fmt.Fprintf(os.Stderr, "direnv: switching\n")
				newEnv, err = loadRC(foundRC, config, oldEnv)
			} else if loadedRC.mtime != foundRC.mtime {
				fmt.Fprintf(os.Stderr, "direnv: reloading\n")
				newEnv, err = loadRC(foundRC, config, oldEnv)
			} else {
				// Nothing to do. Env is loaded and hasn't changed
				return nil
			}
		}

	error:
		if err != nil {
			newEnv = oldEnv
			if foundRC != nil {
				// This should be nearby rc.Load()'s similar statement
				newEnv["DIRENV_DIR"] = "-" + filepath.Dir(foundRC.path)
				newEnv["DIRENV_MTIME"] = fmt.Sprintf("%d", foundRC.mtime)
				newEnv["DIRENV_BACKUP"] = oldEnv.Serialize()
			}
		}

		diff := EnvDiff(env, newEnv)

		diff2 := diff.Filtered()
		if len(diff2) > 0 {
			out := make([]string, len(diff2))
			i := 0
			for key, value := range diff2 {
				if value == "" {
					out[i] = "-" + key
				} else if oldEnv[key] == "" {
					out[i] = "+" + key
				} else {
					out[i] = "~" + key
				}
				i += 1
			}
			sort.Strings(out)
			fmt.Fprintf(os.Stderr, "direnv export: %s\n", strings.Join(out, " "))
		}

		str := EnvToShell(diff, shell)

		fmt.Print(str)
		return

	},
}

func loadRC(rc *RC, config *Config, env Env) (newEnv Env, err error) {
	if !rc.Allowed() {
		return nil, fmt.Errorf("%s is not allowed\n", rc.path)
	}

	argtmpl := `eval "$("%s" stdlib)" >&2 && source_env "%s" >&2 && "%s" dump`
	arg := fmt.Sprintf(argtmpl, config.SelfPath, rc.path, config.SelfPath)
	cmd := exec.Command(config.BashPath, "--noprofile", "--norc", "-c", arg)

	cmd.Stderr = os.Stderr
	cmd.Env = env.ToGoEnv()
	cmd.Dir = filepath.Dir(rc.path)

	out, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("loadRC() failed to run bash: %q", err)
		return
	}

	newEnv, err = ParseEnv(string(out))
	if err != nil {
		err = fmt.Errorf("loadRC() ParseEnv failed: %q", err)
		return
	}

	newEnv["DIRENV_DIR"] = "-" + filepath.Dir(rc.path)
	newEnv["DIRENV_MTIME"] = fmt.Sprintf("%d", rc.mtime)
	newEnv["DIRENV_BACKUP"] = env.Serialize()

	return
}
