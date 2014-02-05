package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type RC struct {
	path      string
	mtime     int64
	allowPath string
}

func FindRC(wd string, allowDir string) *RC {
	rcPath := findUp(wd, ".envrc")
	if rcPath == "" {
		return nil
	}

	return RCFromPath(rcPath, allowDir)
}

func RCFromPath(path string, allowDir string) *RC {
	mtime, err := fileMtime(path)
	if err != nil {
		return nil
	}

	hash, err := fileHash(path)
	if err != nil {
		return nil
	}

	allowPath := filepath.Join(allowDir, hash)
	allowMtime, _ := fileMtime(allowPath)

	if allowMtime > mtime {
		mtime = allowMtime
	}

	return &RC{path, mtime, allowPath}
}

func RCFromEnv(path string, mtime int64) *RC {
	return &RC{path, mtime, ""}
}

func (self *RC) Allow() (err error) {
	if self.allowPath == "" {
		return fmt.Errorf("Cannot allow empty path")
	}
	if err = os.MkdirAll(filepath.Dir(self.allowPath), 0755); err != nil {
		return
	}
	if err = touch(self.allowPath); err != nil {
		return
	}
	self.mtime, err = fileMtime(self.allowPath)
	return
}

func (self *RC) Deny() error {
	return os.Remove(self.allowPath)
}

func (self *RC) Allowed() bool {
	_, err := os.Stat(self.allowPath)
	return err == nil
}

func (self *RC) RelTo(wd string) string {
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

	if !self.Allowed() {
		return nil, fmt.Errorf(NOT_ALLOWED, self.RelTo(wd))
	}

	argtmpl := `eval "$("%s" stdlib)" >&2 && source_env "%s" >&2 && "%s" dump`
	arg := fmt.Sprintf(argtmpl, direnv, self.RelTo(wd), direnv)
	cmd := exec.Command(config.BashPath, "--noprofile", "--norc", "-c", arg)

	cmd.Stderr = os.Stderr
	cmd.Env = env.ToGoEnv()
	cmd.Dir = wd

	out, err := cmd.Output()
	if err != nil {
		return
	}

	newEnv, err = LoadEnv(string(out))
	if err != nil {
		return
	}

	// Save state
	newEnv[DIRENV_DIR] = "-" + filepath.Dir(self.path)
	newEnv[DIRENV_MTIME] = fmt.Sprintf("%d", self.mtime)
	newEnv[DIRENV_DIFF] = env.Diff(newEnv).Serialize()

	return
}

/// Utils

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

func fileMtime(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return fileInfo.ModTime().Unix(), nil
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
	file, err := os.OpenFile(path, os.O_CREATE, 0644)
	if err != nil {
		return
	}
	file.Close()

	t := time.Now()
	return os.Chtimes(path, t, t)
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
