package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
)

// Settings
var libexecDir string

type Command func([]string) error

var commands = map[string]Command{
	// Public commands
	"dump": Dump,
	"diff": Diff,
	// Private commands
	"expand-path": ExpandPath,
	"file-mtime":  FileMtime,
	"file-hash":   FileHash,
	"-h":          Usage,
	"-help":       Usage,
	"--help":      Usage,
}

var publicCommands = []string{"dump", "diff"}

func findCommand(commandName string) (command Command) {
	if command = commands[commandName]; command != nil {
		return
	}
	// First look into the libexec dir
	pathEnv := os.Getenv("PATH")
	os.Setenv("PATH", libexecDir+":"+pathEnv)
	path, err := exec.LookPath("direnv-" + commandName)
	os.Setenv("PATH", pathEnv)

	if err != nil {
		return
	}

	return commandFromPath(path)
}

func commandFromPath(path string) Command {
	return func(args []string) (err error) {
		var files []*os.File
		var process *os.Process
		var procAttr *os.ProcAttr
		var state *os.ProcessState

		files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		procAttr = &os.ProcAttr{"", nil, files, nil}

		// Forward signals to the child process
		c := make(chan os.Signal)
		signal.Notify(c)
		go func() {
			for {
				<-c
				process.Kill()
			}
		}()

		// Execute command
		if process, err = os.StartProcess(path, args, procAttr); err != nil {
			return
		}

		if state, err = process.Wait(); err != nil {
			return
		}

		if !state.Success() {
			return fmt.Errorf("Process failed")
		}

		return nil
	}
}

func availableCommands() []string {
	// TODO: Also return commands available on the path
	return publicCommands
}

func Usage(args []string) error {
	cmds := availableCommands()
	fmt.Fprintf(os.Stderr, "Usage: direnv <%s> [opts]\n", strings.Join(cmds, "|"))
	return nil
}

func init() {
	// Make sure we have the libexec directory available
	libexecDir = os.Getenv("DIRENV_LIBEXEC")
	if libexecDir == "" {
		exePath, err := exec.LookPath(os.Args[0])
		if err != nil {
			exePath = os.Args[0]
		}

		libexecDir, err = filepath.EvalSymlinks(exePath)
		if err != nil {
			libexecDir = exePath
		}

		libexecDir = filepath.Dir(libexecDir)

		// fmt.Fprintf(os.Stderr, "DIRENV_LIBEXEC=%s\n", libexecDir)
		os.Setenv("DIRENV_LIBEXEC", libexecDir)
	}
}

func main() {
	var err error
	var command Command

	// Dispatch
	if len(os.Args) > 1 {
		commandName := os.Args[1]
		args := os.Args[1:]

		command = findCommand(commandName)
		if command == nil {
			fmt.Fprintf(os.Stderr, "Command '%s' not found\n", commandName)
			command = Usage
		}

		err = command(args)
	} else {
		err = Usage(os.Args)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
