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

type EnvDiff map[string]string

func EnvToBash(env EnvDiff) string {
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
			str += "export " + key + "=" + ShellEscape(value) + ";"
		}
	}
	return str
}

func DiffEnv(env1 map[string]string, env2 map[string]string) EnvDiff {
	envDiff := make(EnvDiff)

	for key, _ := range env1 {
		if env2[key] != env1[key] && !IgnoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	// FIXME: I'm sure there is a smarter way to do that
	for key, _ := range env2 {
		if env2[key] != env1[key] && !IgnoredKey(key) {
			envDiff[key] = env2[key]
		}
	}

	return envDiff
}

// A list of keys we don't want to deal with
var IGNORED_KEYS = map[string]bool{"_": true, "PWD": true, "OLDPWD": true, "SHLVL": true, "SHELL": true}

func IgnoredKey(key string) bool {
	if len(key) > 4 && key[:5] == "DIRENV_" {
		return true
	}

	_, found := IGNORED_KEYS[key]
	return found
}

type Env map[string]string

// FIXME: apparently it's possible to have two variables with the
//        same name in the env but I never seen that.
//        If that happens I might have to change the return
//        type signature
func FilteredEnv() Env {
	env := make(Env)

	for _, kv := range os.Environ() {
		kv2 := strings.SplitN(kv, "=", 2)
		// Is there a better way to deconstruct a tuple ?
		key := kv2[0]
		value := kv2[1]

		if !IgnoredKey(key) {
			env[key] = value
		}
	}

	return env
}

func ParseEnv(base64env string) (Env, error) {
	zlibData, err := base64.URLEncoding.DecodeString(base64env)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding: %v", err)
	}

	zlibReader := bytes.NewReader(zlibData)
	w, err := zlib.NewReader(zlibReader)
	if err != nil {
		return nil, fmt.Errorf("zlib opening: %v", err)
	}

	envData := bytes.NewBuffer([]byte{})
	_, err = io.Copy(envData, w)
	if err != nil {
		return nil, fmt.Errorf("zlib decoding: %v", err)
	}
	w.Close()

	env := make(Env)
	err = json.Unmarshal(envData.Bytes(), &env)
	if err != nil {
		return nil, fmt.Errorf("json parsing: %v", err)
	}

	return env, nil
}

func (env Env) Serialize() (string, error) {
	jsonData, err := json.Marshal(env)
	if err != nil {
		return "", err
	}

	zlibData := bytes.NewBuffer([]byte{})
	w := zlib.NewWriter(zlibData)
	w.Write(jsonData)
	w.Close()

	base64Data := base64.URLEncoding.EncodeToString(zlibData.Bytes())

	return base64Data, nil
}
