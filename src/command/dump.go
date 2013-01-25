package command

import (
	"github.com/zimbatm/direnv/src/env"
	"flag"
	"fmt"
)

func Dump(args []string) (err error) {
	flagset := flag.NewFlagSet("direnv dump", flag.ExitOnError)
	flagset.Parse(args[1:])

	e := env.FilteredEnv()
	str, err := e.Serialize()
	if err != nil {
		return
	}
	fmt.Println(str)
	return
}
