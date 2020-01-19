package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RC represents the .envrc file
type RC struct {
	path      string
	allowPath string
	times     FileTimes
	config    *Config
}

// FindRC looks the RC file from the wd, up to the root
func FindRC(wd string, config *Config) *RC {
	rcPath := findUp(wd, ".envrc")
	if rcPath == "" {
		return nil
	}

	return RCFromPath(rcPath, config)
}

// RCFromPath inits the RC from a given path
func RCFromPath(path string, config *Config) *RC {
	hash, err := fileHash(path)
	if err != nil {
		return nil
	}

	allowPath := filepath.Join(config.AllowDir(), hash)

	times := NewFileTimes()

	err = times.Update(path)
	if err != nil {
		return nil
	}

	err = times.Update(allowPath)
	if err != nil {
		return nil
	}

	return &RC{path, allowPath, times, config}
}

// RCFromEnv inits the RC from the environment
func RCFromEnv(path, marshalledTimes string, config *Config) *RC {
	times := NewFileTimes()
	err := times.Unmarshal(marshalledTimes)
	if err != nil {
		return nil
	}
	return &RC{path, "", times, config}
}

// Allow grants the RC as allowed to load
func (rc *RC) Allow() (err error) {
	if rc.allowPath == "" {
		return fmt.Errorf("cannot allow empty path")
	}
	if err = os.MkdirAll(filepath.Dir(rc.allowPath), 0755); err != nil {
		return
	}
	if err = allow(rc.path, rc.allowPath); err != nil {
		return
	}
	err = rc.times.Update(rc.allowPath)
	return
}

// Deny revokes the permission of the RC file to load
func (rc *RC) Deny() error {
	return os.Remove(rc.allowPath)
}

// Allowed checks if the RC file has been granted loading
func (rc *RC) Allowed() bool {
	// happy path is if this envrc has been explicitly allowed, O(1)ish common case
	_, err := os.Stat(rc.allowPath)

	if err == nil {
		return true
	}

	// when whitelisting we want to be (path) absolutely sure we've not been duped with a symlink
	path, err := filepath.Abs(rc.path)
	// seems unlikely that we'd hit this, but have to handle it
	if err != nil {
		return false
	}

	// exact whitelists are O(1)ish to check, so look there first
	if rc.config.WhitelistExact[path] {
		return true
	}

	// finally we check if any of our whitelist prefixes match
	for _, prefix := range rc.config.WhitelistPrefix {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// Path returns the path to the RC file
func (rc *RC) Path() string {
	return rc.path
}

// Touch updates the mtime of the RC file. This is mainly used to trigger a
// reload in direnv.
func (rc *RC) Touch() error {
	return touch(rc.path)
}

const notAllowed = "%s is blocked. Run `direnv allow` to approve its content"

// Load evaluates the RC file and returns the new Env or error.
//
// This functions is key to the implementation of direnv.
func (rc *RC) Load(ctx context.Context, config *Config, env Env) (newEnv Env, err error) {
	wd := config.WorkDir
	direnv := config.SelfPath
	newEnv = env.Copy()
	newEnv[DIRENV_WATCHES] = rc.times.Marshal()
	defer func() {
		rc.RecordState(env, newEnv)
	}()

	if !rc.Allowed() {
		err = fmt.Errorf(notAllowed, rc.Path())
		return
	}

	arg := fmt.Sprintf(
		`eval "$("%s" stdlib)" && __main__ source_env "%s"`,
		direnv,
		rc.Path(),
	)
	cmd := exec.CommandContext(ctx, config.BashPath, "--noprofile", "--norc", "-c", arg)
	cmd.Dir = wd
	cmd.Env = newEnv.ToGoEnv()
	cmd.Stderr = os.Stderr

	if config.DisableStdin {
		cmd.Stdin, err = os.Open(os.DevNull)
		if err != nil {
			return
		}
	} else {
		cmd.Stdin = os.Stdin
	}

	out, err := cmd.Output()
	if err != nil {
		return
	}

	if len(out) > 0 {
		var newEnv2 Env
		newEnv2, err = LoadEnvJSON(out)
		if err != nil {
			return
		}
		if newEnv2["PS1"] != "" {
			logError("PS1 cannot be exported. For more information see https://github.com/direnv/direnv/wiki/PS1")
		}
		newEnv = newEnv2
	}

	return
}

// RecordState just applies the new environment
func (rc *RC) RecordState(env Env, newEnv Env) {
	newEnv[DIRENV_DIR] = "-" + filepath.Dir(rc.path)
	newEnv[DIRENV_DIFF] = env.Diff(newEnv).Serialize()
}

/// Utils

func rootDir(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	i := strings.Index(path[1:], "/")
	if i < 0 {
		return path
	}
	return path[:i+1]
}

func eachDir(path string) (paths []string) {
	path, err := filepath.Abs(path)
	if err != nil {
		return
	}

	paths = []string{path}

	if path == "/" {
		return
	}

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == os.PathSeparator {
			path = path[:i]
			if path == "" {
				path = "/"
			}
			paths = append(paths, path)
		}
	}

	return
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileHash(path string) (hash string, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}

	fd, err := os.Open(path)
	if err != nil {
		return
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte(path + "\n"))
	if err != nil {
		return
	}
	if _, err = io.Copy(hasher, fd); err != nil {
		return
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// Creates a file

func touch(path string) (err error) {
	t := time.Now()
	return os.Chtimes(path, t, t)
}

func allow(path string, allowPath string) (err error) {
	return ioutil.WriteFile(allowPath, []byte(path+"\n"), 0644)
}

func findUp(searchDir string, fileName string) (path string) {
	for _, dir := range eachDir(searchDir) {
		path = filepath.Join(dir, fileName)
		if fileExists(path) {
			return
		}
	}
	return ""
}
