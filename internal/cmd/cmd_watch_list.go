package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// CmdWatchList is `direnv watch-list`
var CmdWatchList = &Cmd{
	Name:    "watch-list",
	Desc:    "Pipe pairs of `mtime path` to stdin to build a list of files to watch.",
	Args:    []string{"[SHELL]"},
	Private: true,
	Action:  actionSimple(watchListCommand),
}

func watchListCommand(env Env, args []string) (err error) {
	var shellName string

	if len(args) >= 2 {
		shellName = args[1]
	} else {
		shellName = "bash"
	}

	shell := DetectShell(shellName)

	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", shellName)
	}

	watches := NewFileTimes()
	watchString, ok := env[DIRENV_WATCHES]
	if ok {
		err = watches.Unmarshal(watchString)
		if err != nil {
			return err
		}
	}

	// Read `mtime path` lines from stdin
	reader := bufio.NewReader(os.Stdin)

	i := 1
	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			elems := strings.SplitN(line, " ", 2)
			if len(elems) != 2 {
				return fmt.Errorf("line %d: expected to contain two elements", i)
			}
			mtime, err := strconv.Atoi(elems[0])
			if err != nil {
				return fmt.Errorf("line %d: %w", i, err)
			}
			path := elems[1][:len(elems[1])-1]

			// add to watches
			err = watches.NewTime(path, int64(mtime), true)
			if err != nil {
				return err
			}
		} else if errors.Is(err, io.EOF) {
			break
		} else {
			return fmt.Errorf("line %d: %w", i, err)
		}
		i++
	}

	e := make(ShellExport)
	e.Add(DIRENV_WATCHES, watches.Marshal())

	os.Stdout.WriteString(shell.Export(e))

	return
}
