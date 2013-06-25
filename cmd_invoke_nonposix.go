// Not a system which supports Exec(). See cmd_invoke_posix for Exec() impl.
// +build !darwin,!freebsd,!linux,!netbsd,!openbsd

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func Invoke(dir string, bash_args []string) (err error) {
	fmt.Println("Non-exec invoke")
	cmd := exec.Command("bash", bash_args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("Failed to invoke shell: %q", err)
	}
	return err
}
