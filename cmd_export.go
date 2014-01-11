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

		loadRC := func() {
			newEnv, err = foundRC.Load(config, oldEnv)
		}

		//fmt.Fprintf(os.Stderr, "%v %v\n", loadedRC, foundRC)

		if loadedRC == nil {
			if foundRC == nil {
				// We're done here.
				return nil
			}

			loadRC()
		} else {
			var backupEnv *EnvDiff
			if backupEnv, err = config.EnvBackup(); err != nil {
				err = fmt.Errorf("EnvBackup() failed: %q", err)
				goto error
			}
			oldEnv = backupEnv.Reverse().Patch(env)
			if foundRC == nil {
				log("unloading")
				newEnv = oldEnv.Filtered()
			} else if loadedRC.path != foundRC.path {
				log("switching")
				loadRC()
			} else if loadedRC.mtime != foundRC.mtime {
				log("reloading")
				loadRC()
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
				newEnv["DIRENV_BACKUP"] = oldEnv.Filtered().Diff(newEnv).Serialize()
			}
		}

		diff := env.Filtered().Diff(newEnv)

		if diff.Any() {
			var out []string
			for key, _ := range diff.Prev {
				_, ok := diff.Next[key]
				if !ok && !ignoredKey(key) {
					out = append(out, "-"+key)
				}
			}
			for key := range diff.Next {
				if ignoredKey(key) {
					continue
				}
				_, ok := diff.Prev[key]
				if ok {
					out = append(out, "~"+key)
				} else {
					out = append(out, "+"+key)
				}
			}
			sort.Strings(out)
			log("export %s", strings.Join(out, " "))
		}

		str := diff.ToShell(shell)
		fmt.Print(str)

		return
	},
}
