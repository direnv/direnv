package cmd

import (
	"fmt"
	"strings"
)

// CmdDiff is `direnv diff`
var CmdDiff = &Cmd{
	Name:    "diff",
	Desc:    "Export list of keys that are different between the current environment and a previous dump",
	Args:    []string{"ENV_DUMP"},
	Private: true,
	Action:  actionSimple(cmdDiffAction),
}

func cmdDiffAction(env Env, args []string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("missing DUMP argument")
	}

	oldDump := args[1]

	oldEnv, err := LoadEnv(oldDump)

	if err != nil {
		return fmt.Errorf("failed to load dump: %w", err)
	}

	diff := env.Diff(oldEnv)

	// If there is no difference, return
	if !diff.Any() {
		return
	}

	// Collect the keys that are different
	var out []string
	for key := range diff.Prev {
		out = append(out, key)
	}
	for key := range diff.Next {
		_, ok := diff.Prev[key]
		if !ok {
			out = append(out, key)
		}
	}

	output := strings.Join(out, ":")
	fmt.Println(output)

	return
}
