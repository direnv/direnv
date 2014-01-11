package main

import (
	"fmt"
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

			newEnv, err = foundRC.Load(config, oldEnv)
		} else {
			var backupEnv Env
			if backupEnv, err = config.EnvBackup(); err != nil {
				err = fmt.Errorf("EnvBackup() failed: %q", err)
				goto error
			}
			oldEnv = backupEnv.Filtered()
			if foundRC == nil {
				log("unloading")
				newEnv = oldEnv
			} else if loadedRC.path != foundRC.path {
				log("switching")
				newEnv, err = foundRC.Load(config, oldEnv)
			} else if loadedRC.mtime != foundRC.mtime {
				log("reloading")
				newEnv, err = foundRC.Load(config, oldEnv)
			} else {
				// Nothing to do. Env is loaded and hasn't changed
				return nil
			}
		}

	error:
		if err != nil {
			newEnv = oldEnv.Filtered()
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
			log("export %s", strings.Join(out, " "))
		}

		str := EnvToShell(diff, shell)

		fmt.Print(str)
		return

	},
}
