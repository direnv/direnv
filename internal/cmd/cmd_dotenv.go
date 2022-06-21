package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/direnv/direnv/v2/pkg/dotenv"
)

// CmdDotEnv is `direnv dotenv [SHELL [PATH_TO_DOTENV]]`
// Transforms a .env file to evaluatable `export KEY=PAIR` statements.
//
// See: https://github.com/bkeepers/dotenv and
//   https://github.com/ddollar/foreman
var CmdDotEnv = &Cmd{
	Name:    "dotenv",
	Desc:    "Transforms a .env file to evaluatable `export KEY=PAIR` statements",
	Args:    []string{"[SHELL]", "[PATH_TO_DOTENV]"},
	Private: true,
	Action:  actionSimple(cmdDotEnvAction),
}

func cmdDotEnvAction(env Env, args []string) (err error) {
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
	if data, err = ioutil.ReadFile(target); err != nil {
		return
	}

	newenv, err = dotenv.Parse(string(data))
	if err != nil {
		return err
	}

	str := newenv.ToShell(shell)
	fmt.Println(str)

	return
}
