package main

import (
	"direnv"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var reverse bool

func init() {
	flag.BoolVar(&reverse, "reverse", false, "Reverses the diff")
	flag.Parse()
}

type EnvDiff map[string]string

// This function and comments have been copied over from Ruby's
// stdlib shellwords.rb library.
func shellEscape(str string) string {
	if str == "" {
		return "''"
	}

	// Treat multibyte characters as is.  It is caller's responsibility
	// to encode the string in the right encoding for the shell
	// environment.
	r := regexp.MustCompile("([^A-Za-z0-9_\\-.,:/@\n])")
	replace := func(match string) string { return "\\\\\\" + match }
	str = r.ReplaceAllStringFunc(str, replace)

	// A LF cannot be escaped with a backslash because a backslash + LF
	// combo is regarded as line continuation and simply ignored.
	str = strings.Replace(str, "\n", "'\n'", -1)
	return str
}

func (env EnvDiff) ToShell() string {
	str := ""
	for key, value := range env {
		// FIXME: This is not exacly as the ruby nil
		if value == "" {
			if key == "PS1" {
				// unsetting PS1 doesn't restore the default in OSX's bash
			} else {
				str += "unset " + key + ";"
			}
		} else {
			str += "export " + key + "=" + shellEscape(value) + ";"
		}
	}
	return str
}

func diffEnv(env1 map[string]string, env2 map[string]string) EnvDiff {
	envDiff := make(EnvDiff)

	for key, _ := range env1 {
		if env2[key] != env1[key] && !direnv.IgnoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	// FIXME: I'm sure there is a smarter way to do that
	for key, _ := range env2 {
		if env2[key] != env1[key] && !direnv.IgnoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	return envDiff
}

func main() {
	oldEnvStr := flag.Arg(0)

	if oldEnvStr == "" {
		log.Fatalln("Missing OLD_ENV argument")
	}

	oldEnv, err := direnv.ParseEnv(oldEnvStr)
	if err != nil {
		log.Fatalln("Parse env error:", err)
	}
	//fmt.Println(oldEnv)

	newEnv := direnv.FilteredEnv()

	var diff EnvDiff
	if reverse {
		diff = diffEnv(oldEnv, newEnv)
	} else {
		diff = diffEnv(newEnv, oldEnv)
	}

	fmt.Print(diff.ToShell())
}
