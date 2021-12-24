package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// CmdWatchDir is `direnv watch-dir SHELL PATH`
var CmdWatchDir = &Cmd{
	Name:    "watch-dir",
	Desc:    "Recursively adds a directory to the list that direnv watches for changes",
	Args:    []string{"SHELL", "DIR"},
	Private: true,
	Action:  actionSimple(watchDirCommand),
}

func watchDirCommand(env Env, args []string) (err error) {
	if len(args) < 3 {
		return fmt.Errorf("a directory is required to add to the list of watches")
	}

	shellName := args[1]
	dir := args[2]

	shell := DetectShell(shellName)

	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", shellName)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("dir '%s' does not exist", dir)
	}

	watches := NewFileTimes()
	watchString, ok := env[DIRENV_WATCHES]
	if ok {
		err = watches.Unmarshal(watchString)
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return watches.NewTime(path, info.ModTime().Unix(), true)
	})
	if err != nil {
		return fmt.Errorf("failed to recursively watch dir '%s': %w", dir, err)
	}

	e := make(ShellExport)
	e.Add(DIRENV_WATCHES, watches.Marshal())

	os.Stdout.WriteString(shell.Export(e))

	return
}
