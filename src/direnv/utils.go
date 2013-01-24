package main

import (
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

// Goes trough all the symlinks and then returns an absolute path
func resolvePath(path string) (retPath string) {
	var err error
	var newPath, wd string

	retPath = path

	for {
		newPath, err = os.Readlink(retPath)
		if err != nil {
			// We reached the end of the chain
			break
		}

		wd = filepath.Dir(retPath)
		retPath = expandPath(newPath, wd)
	}

	return
}
