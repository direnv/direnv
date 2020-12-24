// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js,wasm

package lpenv

import (
	"errors"
)

// ErrNotFound is the error resulting if a path search failed to find an executable file.
var ErrNotFound = errors.New("executable file not found in $PATH")

// LookPathEnv searches for an executable named file in the
// directories named by the PATH environment variable.
// If file contains a slash, it is tried directly and the PATH is not consulted.
// The result may be an absolute path or a path relative to the current directory.
func LookPathEnv(file string, cwd string, env []string) (string, error) {
	// Wasm can not execute processes, so act as if there are no executables at all.
	return "", &Error{file, ErrNotFound}
}
