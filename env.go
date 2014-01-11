package main

import (
	"os"
	"strings"
)

// A list of keys we don't want to deal with
var IGNORED_KEYS = map[string]bool{
	"COMP_WORDBREAKS": true, // Avoids segfaults in bash
	"DIRENV_BASH":     true,
	"DIRENV_CONFIG":   true,
	"OLDPWD":          true,
	"PS1":             true, // Avoids segfaults in bash
	"PWD":             true,
	"SHELL":           true,
	"SHLVL":           true,
	"_":               true,
}

type Env map[string]string

// NOTE:  We don't support having two variables with the same name.
//        I've never seen it used in the wild but accoding to POSIX
//        it's allowed.
func GetEnv() Env {
	env := make(Env)

	for _, kv := range os.Environ() {
		kv2 := strings.SplitN(kv, "=", 2)

		key := kv2[0]
		value := kv2[1]

		env[key] = value
	}

	return env
}

func LoadEnv(base64env string) (env Env, err error) {
	env = make(Env)
	err = unmarshal(base64env, &env)
	return
}

func (env Env) Filtered() Env {
	newEnv := make(Env)

	for key, value := range env {
		if !ignoredKey(key) {
			newEnv[key] = value
		}
	}

	return newEnv
}

func (env Env) ToGoEnv() []string {
	goEnv := make([]string, len(env))
	index := 0
	for key, value := range env {
		goEnv[index] = strings.Join([]string{key, value}, "=")
		index += 1
	}
	return goEnv
}

func (env Env) ToShell(shell Shell) string {
	str := ""

	for key, value := range env {
		str += shell.Export(key, value)
	}

	return str
}

func (env Env) Serialize() string {
	return marshal(env)
}

func (e1 Env) Diff(e2 Env) *EnvDiff {
	return BuildEnvDiff(e1, e2)
}

//// Utils

func ignoredKey(key string) bool {
	if strings.HasPrefix(key, "__fish") {
		return true
	}
	if strings.HasPrefix(key, "DIRENV_") {
		return true
	}
	_, found := IGNORED_KEYS[key]
	return found
}
