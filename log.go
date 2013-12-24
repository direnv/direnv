package main

import (
	"fmt"
	"os"
)

func log(msg string, a ...interface{}) {
	msg = fmt.Sprintf(msg, a...)
	fmt.Fprintf(os.Stderr, "direnv: %s\n", msg)
}
