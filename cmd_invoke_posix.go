// Systems which support Exec(). See cmd_invoke_nonposix for non-Exec() impl.
// +build darwin freebsd linux netbsd openbsd

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Exec into target process.
func Invoke(dir string, bash_args []string) (err error) {
	e := os.Chdir(dir)
	if e != nil {
		return e
	}

	arg0, err := exec.LookPath("bash")
	if err != nil {
		return fmt.Errorf("Can't find bash: %q", err)
	}

	argv := append([]string{arg0}, bash_args...)
	err = syscall.Exec(arg0, argv, os.Environ())
	// Note, remember Exec will never return unless there is an error.
	return err
}
