package main

import (
	_ "embed"
	"github.com/direnv/direnv/v2/internal/cmd"
	"os"
)

var (
	// Configured at compile time
	bashPath string
	//go:embed stdlib.sh
	stdlib string
	//go:embed version.txt
	version string
)

func main() {
	var (
		env  = cmd.GetEnv()
		args = os.Args
	)
	err := cmd.Main(env, args, bashPath, stdlib, version)
	if err != nil {
		os.Exit(1)
	}
}
