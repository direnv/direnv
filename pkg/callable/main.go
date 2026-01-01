// Package callable implements the direnv command-line tool.
package callable

import (
	"context"
	_ "embed"
	"strings"

	"github.com/direnv/direnv/v2/internal/cmd"
)

var (
	// Configured at compile time
	bashPath string
	//go:embed stdlib.sh
	stdlib string
	//go:embed version.txt
	version string
)

func CallableMain(_ context.Context, args []string, env map[string]string) error {

	return cmd.Main(env, args, bashPath, stdlib, strings.TrimSpace(version))
}
