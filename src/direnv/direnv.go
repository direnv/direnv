package direnv

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// A list of keys we don't want to deal with
var IGNORED_KEYS []string

func init() {
	IGNORED_KEYS = []string{"_", "PWD", "OLDPWD", "SHLVL", "SHELL", "DIRENV_BACKUP", "DIRENV_LIBEXEC"}
}

func IgnoredKey(key string) bool {
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
//        same name in the env but I never saw that.
//        If that happens I might to have to change the return
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
	zlibData, err := base64.StdEncoding.DecodeString(base64env)
	if err != nil {
		return nil, fmt.Errorf("base64 decoding: %v", err)
	}

	zlibReader := bytes.NewReader(zlibData)
	w, err := zlib.NewReader(zlibReader)
	if err != nil {
		return nil, fmt.Errorf("zlib opening: %v", err)
	}

	envData := new(bytesReceiver)
	io.Copy(envData, w)
	//_, err = io.Copy(w, envData)
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

	var zlibData bytes.Buffer
	w := zlib.NewWriter(&zlibData)
	w.Write(jsonData)
	w.Close()

	base64Data := base64.StdEncoding.EncodeToString(zlibData.Bytes())

	return base64Data, nil
}

// Erk. There's probably something in the stdlib that does what this
// micro type does.
type bytesReceiver struct {
	bytes []byte
}

func (self *bytesReceiver) Bytes() []byte {
	return self.bytes
}

func (self *bytesReceiver) Write(data []byte) (int, error) {
	// HACK: We know only one write is going to happen :)
	// ..... but still make sure the assertion is correct
	if len(self.bytes) > 0 {
		return 0, errors.New("You fool, thinking you can be right")
	}
	self.bytes = data
	return len(data), nil
}
