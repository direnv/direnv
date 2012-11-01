package main

import (
	"fmt"
	"flag"
	"env"
)

func Diff(args []string) (err error) {
	var reverse bool

	flagset := flag.NewFlagSet("direnv dump", flag.ExitOnError)
	flagset.BoolVar(&reverse, "reverse", false, "Reverses the diff")
	err = flagset.Parse(args)
	if err != nil {
		return
	}

	oldEnvStr := flagset.Arg(0)

	if oldEnvStr == "" {
		return fmt.Errorf("Missing OLD_ENV argument")
	}

	oldEnv, err := env.ParseEnv(oldEnvStr)
	if err != nil {
		return fmt.Errorf("Parse env error: %v", err)
	}

	newEnv := env.FilteredEnv()

	var diff env.EnvDiff
	if reverse {
		diff = env.Diff(oldEnv, newEnv)
	} else {
		diff = env.Diff(newEnv, oldEnv)
	}

	fmt.Print(diff.ToShell())

	return
}
