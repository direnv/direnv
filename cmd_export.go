package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// `direnv export $0`
var CmdExport = &Cmd{
	Name:    "export",
	Desc:    "loads an .envrc and prints the diff in terms of exports",
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
			fmt.Fprintf(os.Stderr, "direnv: %s\n", strings.Join(out, ","))
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

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	r2 := bufio.NewReader(r)

	attr := &os.ProcAttr{
		Dir:   filepath.Dir(rc.path),
		Env:   env.ToGoEnv(),
		Files: []*os.File{os.Stdin, w, os.Stderr},
	}

	command := fmt.Sprintf(`eval "$("%s" stdlib)" >&2 && source_env "%s" >&2 && "%s" dump`, config.SelfPath, rc.path, config.SelfPath)

	process, err := os.StartProcess(config.BashPath, []string{"bash", "-c", command}, attr)
	if err != nil {
		return nil, err
	}

	output, err := r2.ReadString('\n')
	if err != nil {
		panic(err)
	}

	_, err = process.Wait()
	if err != nil {
		return nil, err
	}

	newEnv, err = ParseEnv(output)
	if err != nil {
		return
	}

	newEnv["DIRENV_DIR"] = "-" + filepath.Dir(rc.path)
	newEnv["DIRENV_MTIME"] = fmt.Sprintf("%d", rc.mtime)
	newEnv["DIRENV_BACKUP"] = env.Serialize()

	return newEnv, nil
}
