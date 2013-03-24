//
// Commands that we want to expose in the stdlib.
// Generally they exist because of cross-platform issues.
//

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Load(env Env, args []string) (err error) {
	flagset := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagset.Parse(args[1:])

	str := flagset.Arg(0)
	env, err = ParseEnv(str)
	fmt.Println(env)
	return
}

func Dump(env Env, args []string) (err error) {
	flagset := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagset.Parse(args[1:])

	fmt.Println(env.Filtered().Serialize())
	return
}

func Export(env Env, args []string) (err error) {
	var newEnv Env
	var loadedRC *RC
	var foundRC *RC
	var config *Config
	var target string

	if len(args) > 0 {
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
	foundRC = config.FoundRC()

	if loadedRC != nil {
		var backupEnv Env
		if backupEnv, err = config.EnvBackup(); err != nil {
			return
		}

		if foundRC == nil {
			fmt.Fprintf(os.Stderr, "Unloading %s\n", loadedRC.path)
			newEnv = backupEnv
		} else if loadedRC.path != foundRC.path {
			fmt.Fprintf(os.Stderr, "Switching from %s to %s\n", loadedRC.path, foundRC.path)
			newEnv, err = foundRC.Load(backupEnv, config.ExecDir)
		} else if loadedRC.mtime != foundRC.mtime {
			fmt.Fprintf(os.Stderr, "Reloading %s\n", loadedRC.path)
			newEnv, err = foundRC.Load(backupEnv, config.ExecDir)
		} else {
			// Nothing to do. Env is loaded and hasn't changed
			return nil
		}
	} else {
		if foundRC == nil {
			// Done here
			return
		}

		fmt.Fprintf(os.Stderr, "Loading %s\n", foundRC.path)
		newEnv, err = foundRC.Load(env, config.ExecDir)
	}

	if err != nil {
		return
	}

	diff := EnvDiff(env, newEnv)

	diff2 := diff.Filtered()
	if len(diff2) > 0 {
		fmt.Fprintf(os.Stderr, "Changed: %s\n", strings.Join(stringKeys(diff2), ","))
	}

	str := EnvToShell(diff, shell)

	fmt.Println(str)
	return
}

func expandPath(path, relTo string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(relTo, path))
}

func ExpandPath(env Env, args []string) (err error) {
	var path string

	flagset := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagset.Parse(args[1:])

	path = flagset.Arg(0)
	if path == "" {
		return fmt.Errorf("PATH missing")
	}

	if !filepath.IsAbs(path) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		relTo := flagset.Arg(1)
		if relTo == "" {
			relTo = wd
		} else {
			relTo = expandPath(relTo, wd)
		}

		path = expandPath(path, relTo)
	}

	_, err = fmt.Println(path)

	return
}

// Utils

func stringKeys(hash map[string]string) []string {
	keys := make([]string, len(hash))
	i := 0
	for key, _ := range hash {
		keys[i] = key
		i += 1
	}
	return keys
}
