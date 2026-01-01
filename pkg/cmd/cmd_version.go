package cmd

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

// CmdVersion is `direnv version`
var CmdVersion = &Cmd{
	Name:    "version",
	Desc:    "prints the version or checks that direnv is older than VERSION_AT_LEAST.",
	Args:    []string{"[VERSION_AT_LEAST]"},
	Aliases: []string{"--version"},
	Action: actionSimple(func(_ Env, args []string) error {
		semVersion := ensureVPrefixed(version)
		if len(args) > 1 {
			atLeast := ensureVPrefixed(args[1])
			if !semver.IsValid(atLeast) {
				return fmt.Errorf("%s is not a valid semver version", atLeast)
			}
			cmp := semver.Compare(semVersion, atLeast)
			if cmp < 0 {
				return fmt.Errorf("current version %s is older than the desired version %s", semVersion, atLeast)
			}
		} else {
			fmt.Println(version)
		}
		return nil
	}),
}

func ensureVPrefixed(version string) string {
	if !strings.HasPrefix(version, "v") {
		return "v" + version
	}
	return version
}
