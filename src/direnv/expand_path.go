package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Makes a path absolute, relative to another "wd" path
func expandPath(path string, wd string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(wd, path)
}

func ExpandPath(args []string) (err error) {
	var path, wd string

	flagset := flag.NewFlagSet("direnv expand-path", flag.ExitOnError)
	flagset.Parse(args[1:])

	path = flagset.Arg(0)
	if path == "" {
		return fmt.Errorf("PATH missing")
	}

	wd = flagset.Arg(1)
	if wd == "" {
		if wd, err = os.Getwd(); err != nil {
			return
		}
	}

	absPath := expandPath(path, wd)
	fmt.Println(absPath)

	return nil
}
