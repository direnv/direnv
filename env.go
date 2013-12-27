package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Env map[string]string

func EnvToShell(env Env, shell Shell) string {
	str := ""
	for key, value := range env {
		// FIXME: This is not exacly as the ruby nil
		if value == "" {
			if key == "PS1" {
				// unsetting PS1 doesn't restore the default in OSX's bash
			} else {
				str += shell.Unset(key)
			}
		} else {
			str += shell.Export(key, value)
		}
	}
	return str
}

func EnvDiff(env1 map[string]string, env2 map[string]string) Env {
	envDiff := make(Env)

	for key := range env1 {
		if env2[key] != env1[key] && !ignoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	// FIXME: I'm sure there is a smarter way to do that
	for key := range env2 {
		if env2[key] != env1[key] && !ignoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	return envDiff
}

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

func ignoredKey(key string) bool {
	if len(key) > 6 && key[0:6] == "__fish" {
		return true
	}
	_, found := IGNORED_KEYS[key]
	return found
}

func direnvVar(key string) bool {
	return strings.HasPrefix(key, "DIRENV_")
}

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

func (env Env) Filtered() Env {
	newEnv := make(Env)

	for key, value := range env {
		if !ignoredKey(key) && !direnvVar(key) {
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

func ParseEnv(base64env string) (Env, error) {
	zlibData, err := base64.URLEncoding.DecodeString(base64env)
	if err != nil {
		return nil, fmt.Errorf("ParseEnv() base64 decoding: %v", err)
	}

	zlibReader := bytes.NewReader(zlibData)
	w, err := zlib.NewReader(zlibReader)
	if err != nil {
		return nil, fmt.Errorf("ParseEnv() zlib opening: %v", err)
	}

	envData := bytes.NewBuffer([]byte{})
	_, err = io.Copy(envData, w)
	if err != nil {
		return nil, fmt.Errorf("ParseEnv() zlib decoding: %v", err)
	}
	w.Close()

	env := make(Env)
	err = json.Unmarshal(envData.Bytes(), &env)
	if err != nil {
		return nil, fmt.Errorf("ParseEnv() json parsing: %v", err)
	}

	return env, nil
}

func (env Env) Serialize() string {
	// We can safely ignore the err because it's only thrown
	// for unsupported datatype. We know that a map[string]string
	// is supported.
	jsonData, err := json.Marshal(env)
	if err != nil {
		panic(fmt.Errorf("Serialize(): %q", err))
	}

	zlibData := bytes.NewBuffer([]byte{})
	w := zlib.NewWriter(zlibData)
	w.Write(jsonData)
	w.Close()

	base64Data := base64.URLEncoding.EncodeToString(zlibData.Bytes())

	return base64Data
}
