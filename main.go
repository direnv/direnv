// Package main implements the direnv command-line tool.
package main

import (
	"context"
	"os"

	"github.com/yaklabco/direnv/v2/pkg/callable"
	"github.com/yaklabco/direnv/v2/pkg/cmd"
)

func main() {
	if err := callable.CallableMain(context.Background(), os.Args, cmd.GetEnv()); err != nil {
		os.Exit(1)
	}
}
