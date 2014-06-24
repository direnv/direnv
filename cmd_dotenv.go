package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var DOTENV_REG = regexp.MustCompile("(?:export\\s+)?([\\w\\.]+)(?:\\s*=\\s*|:\\s+?)(.*)")
var DOTENV_LF_REG = regexp.MustCompile("\\\\n")
var DOTENV_ESC_REG = regexp.MustCompile("\\\\.")

func ParseDotEnv(data string) Env {
	var dotenv = make(Env)

	result := DOTENV_REG.FindAllStringSubmatch(data, -1)
	for _, match := range result {
		key := match[1]
		value := strings.TrimSpace(match[2])

		if value[0:1] == "'" && value[len(value)-1:] == "'" {
			value = value[1 : len(value)-1]
		} else if value[0:1] == `"` && value[len(value)-1:] == `"` {
			value = value[1 : len(value)-1]
			value = DOTENV_LF_REG.ReplaceAllString(value, "\n")
			value = DOTENV_ESC_REG.ReplaceAllStringFunc(value, func(str string) string {
				return str[1:2]
			})
		}

		dotenv[key] = value
	}

	return dotenv
}

// `direnv private dotenv [SHELL [PATH_TO_DOTENV]]`
// Transforms a .env file to evaluatable `export KEY=PAIR` statements.
//
// See: https://github.com/bkeepers/dotenv and
//   https://github.com/ddollar/foreman
var CmdDotEnv = &Cmd{
	Name:    "dotenv",
	Desc:    "Transforms a .env file to evaluatable `export KEY=PAIR` statements",
	Args:    []string{"[SHELL]", "[PATH_TO_DOTENV]"},
	Private: true,
	Fn: func(env Env, args []string) (err error) {
		var shell Shell
		var target string

		if len(args) > 1 {
			shell = DetectShell(args[1])
		} else {
			shell = BASH
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

		env = ParseDotEnv(string(data))
		str := env.ToShell(shell)
		fmt.Println(str)

		return
	},
}
