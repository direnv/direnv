package main

import (
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

type RC struct {
	path      string
	allowPath string
	times     FileTimes
}

func FindRC(wd string, allowDir string) *RC {
	rcPath := findUp(wd, ".envrc")
	if rcPath == "" {
		return nil
	}

	return RCFromPath(rcPath, allowDir)
}

func RCFromPath(path string, allowDir string) *RC {
	hash, err := fileHash(path)
	if err != nil {
		return nil
	}

	allowPath := filepath.Join(allowDir, hash)

	times := NewFileTimes()
	times.Update(path)
	times.Update(allowPath)

	return &RC{path, allowPath, times}
}

func RCFromEnv(path, marshalled_times string) *RC {
	times := NewFileTimes()
	times.Unmarshal(marshalled_times)
	return &RC{path, "", times}
}

func (self *RC) Allow() (err error) {
	if self.allowPath == "" {
		return fmt.Errorf("Cannot allow empty path")
	}
	if err = os.MkdirAll(filepath.Dir(self.allowPath), 0755); err != nil {
		return
	}
	if err = allow(self.path, self.allowPath); err != nil {
		return
	}
	self.times.Update(self.allowPath)
	return
}

func (self *RC) Deny() error {
	return os.Remove(self.allowPath)
}

func (self *RC) Allowed() bool {
	_, err := os.Stat(self.allowPath)
	return err == nil
}

// Makes the path relative to the current directory. Except when both paths
// are completely different.
// Eg:  /home/foo and /home/bar => ../foo
// But: /home/foo and /tmp/bar  => /home/foo
func (self *RC) RelTo(wd string) string {
	if rootDir(wd) != rootDir(self.path) {
		return self.path
	}
	x, err := filepath.Rel(wd, self.path)
	if err != nil {
		panic(err)
	}
	return x
}

func (self *RC) Touch() error {
	return touch(self.path)
}

const NOT_ALLOWED = "%s is blocked. Run `direnv allow` to approve its content."

func (self *RC) Load(config *Config, env Env) (newEnv Env, err error) {
	wd := config.WorkDir
	direnv := config.SelfPath
	shellEnv := env.Copy()
	shellEnv[DIRENV_WATCHES] = self.times.Marshal()

	if !self.Allowed() {
		return nil, fmt.Errorf(NOT_ALLOWED, self.RelTo(wd))
	}

	argtmpl := `eval "$("%s" stdlib)" >&2 && source_env "%s" >&2 && "%s" dump`
	arg := fmt.Sprintf(argtmpl, direnv, self.RelTo(wd), direnv)
	cmd := exec.Command(config.BashPath, "--noprofile", "--norc", "-c", arg)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = shellEnv.ToGoEnv()
	cmd.Dir = wd

	out, err := cmd.Output()
	if err != nil {
		return
	}

	newEnv, err = LoadEnv(string(out))
	if err != nil {
		return
	}

	self.RecordState(env, newEnv)

	return
}

func (self *RC) RecordState(env Env, newEnv Env) {
	newEnv[DIRENV_DIR] = "-" + filepath.Dir(self.path)
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
	hasher.Write([]byte(path + "\n"))
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
