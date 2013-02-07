package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	// "io/ioutil"
	"os"
	"path/filepath"
)

type RC struct {
	path      string
	mtime     int64
	hash      string
	allowPath string
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

func fileMtime(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return fileInfo.ModTime().Unix(), nil
}

func fileHash(path string) (string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	io.Copy(hasher, fd)
	num := hasher.Sum(nil)

	return fmt.Sprintf("%x", num), nil
}

// Creates a file
func touch(path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	//	fd.Write([]byte{})
	fd.Close()
	return nil
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

func FindRC(wd string, allowDir string) *RC {
	if wd == "" {
		var err error
		if wd, err = os.Getwd(); err != nil {
			return nil
		}
	}

	rcPath := findUp(wd, ".envrc")
	if rcPath == "" {
		return nil
	}

	return LoadRC(rcPath, allowDir)
}

func LoadRC(path string, allowDir string) *RC {
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

	return &RC{path, mtime, hash, allowPath}
}

func RCFromEnv(path string, mtime int64, hash string, allowDir string) *RC {
	allowPath := filepath.Join(allowDir, hash)
	return &RC{path, mtime, hash, allowPath}
}

func (self *RC) Allow() (err error) {
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

func (rc *RC) Load(env Env, workDir string) (newEnv Env, err error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	r2 := bufio.NewReader(r)

	attr := &os.ProcAttr{
		Dir:   filepath.Dir(rc.path),
		Env:   env.ToGoEnv(),
		Files: []*os.File{os.Stdin, w, os.Stderr},
	}

	command := fmt.Sprintf("eval `%s/direnv private stdlib` >&2; source %s >&2 && %s/direnv private dump", workDir, rc.path, workDir)

	process, err := os.StartProcess("/bin/bash", []string{"bash", "-c", command}, attr)
	if err != nil {
		return nil, err
	}

	output, err := r2.ReadString('\n')
	if err != nil {
		panic(err)
	}

	_, err = process.Wait()
	if err != nil {
		return nil, err
	}

	newEnv, err = ParseEnv(output)
	if err != nil {
		return
	}

	newEnv["DIRENV_DIR"] = filepath.Dir(rc.path)
	newEnv["DIRENV_MTIME"] = fmt.Sprintf("%d", rc.mtime)
	newEnv["DIRENV_BACKUP"] = env.Serialize()

	return newEnv, nil
}
