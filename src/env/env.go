package env

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

// A list of keys we don't want to deal with
var IGNORED_KEYS = []string{"_", "PWD", "OLDPWD", "SHLVL", "SHELL"}

func IgnoredKey(key string) bool {
	if len(key) > 4 && key[:5] == "DIRENV_" {
		return true
	}
	// FIXME: Is there a higher-level function for that ? Eg. indexOf in JavaScript
	for _, ikey := range IGNORED_KEYS {
		if ikey == key {
			return true
		}
	}
	return false
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
