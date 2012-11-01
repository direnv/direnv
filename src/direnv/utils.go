package main

import (
	"os"
	"path/filepath"
)

// Goes trough all the symlinks and then returns an absolute path
func resolvePath(path string) (retPath string) {
	retPath = path
	for {
		newPath, err := os.Readlink(retPath)
		if err != nil {
			break
		}
		retPath = newPath
	}
	absPath, err := filepath.Abs(retPath)
	if err == nil {
		retPath = absPath
	}

	return
}
