package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
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
	"-h": Usage,
	"-help": Usage,
	"--help": Usage,
}

var publicCommands = []string{ "dump", "diff" }

func availableCommands() []string {
	// TODO: Also return commands available on the path
	return publicCommands;
}

func Usage(args []string) error {
	cmds := availableCommands();
	fmt.Printf("Usage: direnv <%s> [opts]\n", strings.Join(cmds, "|"))
	return nil
}

func init() {
	// Make sure we have the libexec directory available
	libexecDir = os.Getenv("DIRENV_LIBEXEC")
	if libexecDir == "" {
		libexecDir = path.Dir(resolvePath(os.Args[0]))
		fmt.Println(libexecDir)
		os.Setenv("DIRENV_LIBEXEC", libexecDir)
	}
}

func main() {
	var err error

	// Dispatch
	if len(os.Args) > 1 {
		commandName := os.Args[1]
		command := commands[commandName]
		args := os.Args[1:]

		if command != nil {
			err = command(args)
		} else {
			// First look into the libexec dir
			pathEnv := os.Getenv("PATH")
			os.Setenv("PATH", libexecDir + ":" + pathEnv)
			path, err := exec.LookPath("direnv-" + commandName)
			os.Setenv("PATH", pathEnv)

			if err != nil {
				fmt.Println("Command not found")
				err = Usage(args)
			} else {
				var files []*os.File
				files = append(files, os.Stdin)
				files = append(files, os.Stdout)
				files = append(files, os.Stderr)

				// Execute command
				// TODO: Forward signal
				process, err := os.StartProcess(path, args, &os.ProcAttr{"", nil, files, nil})
				if err != nil {
					process.Wait()
				}
			}
		}
	} else {
		err = Usage(os.Args)
	}

	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}
