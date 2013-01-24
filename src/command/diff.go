package command

import (
	"../env"
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

	fmt.Println(env.ToBash(diff))

	return
}
