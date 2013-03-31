package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// This is run by the shell before every prompt
func Export(env Env, args []string) (err error) {
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

		fmt.Fprintf(os.Stderr, "Loading %s\n", foundRC.path)
		newEnv, err = foundRC.Load(config, oldEnv)
	} else {
		var backupEnv Env
		if backupEnv, err = config.EnvBackup(); err != nil {
			goto error
		}
		oldEnv = backupEnv.Filtered()
		if foundRC == nil {
			fmt.Fprintf(os.Stderr, "Unloading %s\n", loadedRC.path)
			newEnv = oldEnv
		} else if loadedRC.path != foundRC.path {
			fmt.Fprintf(os.Stderr, "Switching from %s to %s\n", loadedRC.path, foundRC.path)
			newEnv, err = foundRC.Load(config, oldEnv)
		} else if loadedRC.mtime != foundRC.mtime {
			fmt.Fprintf(os.Stderr, "Reloading %s\n", loadedRC.path)
			newEnv, err = foundRC.Load(config, oldEnv)
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
		fmt.Fprintf(os.Stderr, "Changed: %s\n", strings.Join(mapKeys(diff2), ","))
	}

	str := EnvToShell(diff, shell)

	fmt.Println(str)
	return

}

func mapKeys(hash map[string]string) []string {
	keys := make([]string, len(hash))
	i := 0
	for key, _ := range hash {
		keys[i] = key
		i += 1
	}
	return keys
}
