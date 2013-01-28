package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type RC struct {
	Path  string
	Mtime int64
	Hash  string
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

func findUp(searchDir string, fileName string) (path string) {
	for _, dir := range eachDir(searchDir) {
		path = filepath.Join(dir, fileName)
		if fileExists(path) {
			return
		}
	}
	return ""
}

func FindRC(wd string) *RC {
	rcPath := findUp(wd, ".envrc")
	if rcPath == "" {
		return nil
	}

	mtime, err := fileMtime(rcPath)
	if err != nil {
		return nil
	}
	hash, err := fileHash(rcPath)
	if err != nil {
		return nil
	}

	return &RC{rcPath, mtime, hash}
}
