// Package callable implements the direnv command-line tool.
package callable

import (
	"context"
	_ "embed"
	"strings"

	"github.com/yaklabco/direnv/v2/pkg/cmd"
)

var (
	// Configured at compile time
	bashPath string
	//go:embed stdlib.sh
	stdlib string
	//go:embed version.txt
	version string
)

// CallableMain is the entrypoint for the direnv command-line tool.
// It processes the provided arguments and environment variables,
// then delegates to the internal cmd.Main function with necessary parameters.
func CallableMain(_ context.Context, args []string, env map[string]string) error { //nolint:revive // We're okay with the stutter here.
	return cmd.Main(env, args, bashPath, stdlib, strings.TrimSpace(version))
}
