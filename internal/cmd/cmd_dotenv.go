package cmd

import (
	"path/filepath"
	"fmt"
	"os"
	"github.com/direnv/direnv/v2/pkg/dotenv"
)

// CmdDotEnv is `direnv dotenv [SHELL [PATH_TO_DOTENV]]`
// Transforms a .env file to evaluatable `export KEY=PAIR` statements.
//
// See: https://github.com/bkeepers/dotenv and https://github.com/ddollar/foreman
var CmdDotEnv = &Cmd{
	Name:    "dotenv",
	Desc:    "Transforms a .env file to evaluatable `export KEY=PAIR` statements",
	Args:    []string{"[SHELL]", "[PATH_TO_DOTENV]"},
	Private: true,
	Action:  actionSimple(cmdDotEnvAction),
}

func cmdDotEnvAction(_ Env, args []string) (err error) {
	var shell Shell
	var newenv Env
	var target string

	if len(args) > 1 {
		shell = DetectShell(args[1])
	} else {
		shell = Bash
	}

	if len(args) > 2 {
		target = args[2]
	}

	if target == "" {
		target = ".env"
	}

	var data []byte
	if data, err = os.ReadFile(target); err != nil {
		return
	}
	// Set PWD env var to the directory the .env file resides in This results in
	// the least amount of surprise, as a dotenv file is most often defined
	// in the same directory it's loaded from, so refering to PWD should
	// match the directory of the dotenv file.
	path, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	pwd := filepath.Dir(path)
	os.Setenv("PWD", pwd)
	newenv, err = dotenv.Parse(string(data))
	if err != nil {
		return err
	}

	str := newenv.ToShell(shell)
	fmt.Println(str)

	return
}
