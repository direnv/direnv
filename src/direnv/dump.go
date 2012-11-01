package main

import (
	"fmt"
	"flag"
	"env"
)

func Dump(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv dump", flag.ExitOnError)
	err = flagset.Parse(args)
	if err != nil {
		return
	}

	e := env.FilteredEnv()
	str, err := e.Serialize()
	if err != nil {
		return
	}
	fmt.Print(str)
	return
}

