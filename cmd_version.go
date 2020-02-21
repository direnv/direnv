package main

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

// CmdVersion is `direnv version`
var CmdVersion = &Cmd{
	Name:    "version",
	Desc:    "prints the version (" + Version + ") or checks that direnv is older than VERSION_AT_LEAST.",
	Args:    []string{"[VERSION_AT_LEAST]"},
	Aliases: []string{"--version"},
	Action: actionSimple(func(env Env, args []string) error {
		semVersion := "v" + Version
		if len(args) > 1 {
			atLeast := args[1]
			if !strings.HasPrefix(atLeast, "v") {
				atLeast = "v" + atLeast
			}
			if !semver.IsValid(atLeast) {
				return fmt.Errorf("%s is not a valid semver version", atLeast)
			}
			cmp := semver.Compare(semVersion, atLeast)
			if cmp < 0 {
				return fmt.Errorf("current version %s is older than the desired version %s", semVersion, atLeast)
			}
		} else {
			fmt.Println(Version)
		}
		return nil
	}),
}
