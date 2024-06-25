package cmd

import (
	"errors"
	"os"
	"strings"
)

// lookPathUnix is similar to os/exec.LookPath except we pass in the PATH
//
// Also see:
// - https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/os/exec/lp_unix.go;l=52
// - https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/os/exec/lp_unix_test.go
func lookPathUnix(file string, pathenv string) (string, error) {
	if strings.Contains(file, "/") {
		err := findExecutable(file)
		if err == nil {
			return file, nil
		}
		return "", err
	}
	if pathenv == "" {
		return "", errNotFound
	}
	for _, dir := range strings.Split(pathenv, ":") {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := dir + "/" + file
		if err := findExecutable(path); err == nil {
			return path, nil
		}
	}
	return "", errNotFound
}

// Similar to os/exec.LookPath except we pass in the PATH
func lookPath(file string, pathenv string) (string, error) {
	if Cygpath != nil {
		return Cygpath.LookPath(file, pathenv)
	}
	return lookPathUnix(file, pathenv)
}

// ErrNotFound is the error resulting if a path search failed to find an executable file.
var errNotFound = errors.New("executable file not found in $PATH")

func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}
