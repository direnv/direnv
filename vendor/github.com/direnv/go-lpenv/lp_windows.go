// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lpenv

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrNotFound is the error resulting if a path search failed to find an executable file.
var ErrNotFound = errors.New("executable file not found in %PATH%")

func chkStat(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if d.IsDir() {
		return os.ErrPermission
	}
	return nil
}

func hasExt(file string) bool {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		return false
	}
	return strings.LastIndexAny(file, `:\/`) < i
}

func findExecutable(cwd string, file string, exts []string) (string, error) {
	file = JoinRel(cwd, file)
	if len(exts) == 0 {
		return file, chkStat(file)
	}
	if hasExt(file) {
		if chkStat(file) == nil {
			return file, nil
		}
	}
	for _, e := range exts {
		if f := file + e; chkStat(f) == nil {
			return f, nil
		}
	}
	return "", os.ErrNotExist
}

// LookPathEnv searches for an executable named file in the
// directories named by the PATH environment variable.
// If file contains a slash, it is tried directly and the PATH is not consulted.
// LookPathEnv also uses PATHEXT environment variable to match
// a suitable candidate.
// The result may be an absolute path or a path relative to the current directory.
func LookPathEnv(file string, cwd string, env []string) (string, error) {
	var exts []string
	x := Getenv(env, `PATHEXT`)
	if x != "" {
		for _, e := range strings.Split(strings.ToLower(x), `;`) {
			if e == "" {
				continue
			}
			if e[0] != '.' {
				e = "." + e
			}
			exts = append(exts, e)
		}
	} else {
		exts = []string{".com", ".exe", ".bat", ".cmd"}
	}

	if strings.ContainsAny(file, `:\/`) {
		if f, err := findExecutable(cwd, file, exts); err == nil {
			return f, nil
		} else {
			return "", &Error{file, err}
		}
	}
	if f, err := findExecutable(cwd, file, exts); err == nil {
		return f, nil
	}
	path := Getenv(env, "path")
	for _, dir := range filepath.SplitList(path) {
		if f, err := findExecutable(cwd, filepath.Join(dir, file), exts); err == nil {
			return f, nil
		}
	}
	return "", &Error{file, ErrNotFound}
}
