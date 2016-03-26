package main

import (
	"os"
	"strings"
)

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

func (env Env) CleanContext() {
	delete(env, DIRENV_DIR)
	delete(env, DIRENV_WATCHES)
	delete(env, DIRENV_DIFF)
}

func LoadEnv(base64env string) (env Env, err error) {
	env = make(Env)
	err = unmarshal(base64env, &env)
	return
}

func (env Env) Copy() Env {
	newEnv := make(Env)

	for key, value := range env {
		newEnv[key] = value
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
	e := make(ShellExport)

	for key, value := range env {
		e.Add(key, value)
	}

	return shell.Export(e)
}

func (env Env) Serialize() string {
	return marshal(env)
}

func (e1 Env) Diff(e2 Env) *EnvDiff {
	return BuildEnvDiff(e1, e2)
}
