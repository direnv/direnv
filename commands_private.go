//
// Commands that we want to expose in the stdlib.
// Generally they exist because of cross-platform issues.
//

package main

import (
	"flag"
	"fmt"
)

func Diff(args []string) (err error) {
	var reverse bool

	flagset := flag.NewFlagSet("direnv dump", flag.ExitOnError)
	flagset.BoolVar(&reverse, "reverse", false, "Reverses the diff")
	flagset.Parse(args[1:])

	oldEnvStr := flagset.Arg(0)

	if oldEnvStr == "" {
		return fmt.Errorf("Missing OLD_ENV argument")
	}

	oldEnv, err := ParseEnv(oldEnvStr)
	if err != nil {
		return fmt.Errorf("Parse env error: %v", err)
	}

	newEnv := GetEnv().Filtered()

	var diff Env
	if reverse {
		diff = EnvDiff(oldEnv, newEnv)
	} else {
		diff = EnvDiff(newEnv, oldEnv)
	}

	fmt.Println(EnvToBash(diff))

	return
}

func Dump(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv dump", flag.ExitOnError)
	flagset.Parse(args[1:])

	e := GetEnv().Filtered()
	str, err := e.Serialize()
	if err != nil {
		return
	}
	fmt.Println(str)
	return
}

func Export(args []string) (err error) {
	// TODO

	return
}
