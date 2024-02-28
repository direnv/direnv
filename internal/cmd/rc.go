package cmd

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

// RC represents the .envrc or .env file
type RC struct {
	path      string
	allowPath string
	denyPath  string
	times     FileTimes
	config    *Config
}

// FindRC looks for ".envrc" and ".env" files up in the file hierarchy.
func FindRC(wd string, config *Config) (*RC, error) {
	rcPath := findEnvUp(wd, config.LoadDotenv)
	if rcPath == "" {
		return nil, nil
	}

	return RCFromPath(rcPath, config)
}

// RCFromPath inits the RC from a given path
func RCFromPath(path string, config *Config) (*RC, error) {
	fileHash, err := fileHash(path)
	if err != nil {
		return nil, err
	}

	allowPath := filepath.Join(config.AllowDir(), fileHash)

	pathHash, err := pathHash(path)
	if err != nil {
		return nil, err
	}

	denyPath := filepath.Join(config.DenyDir(), pathHash)

	times := NewFileTimes()

	err = times.Update(path)
	if err != nil {
		return nil, err
	}

	err = times.Update(allowPath)
	if err != nil {
		return nil, err
	}

	err = times.Update(denyPath)
	if err != nil {
		return nil, err
	}

	return &RC{path, allowPath, denyPath, times, config}, nil
}

// RCFromEnv inits the RC from the environment
func RCFromEnv(path, marshalledTimes string, config *Config) *RC {
	fileHash, err := fileHash(path)
	if err != nil {
		return nil
	}

	allowPath := filepath.Join(config.AllowDir(), fileHash)

	times := NewFileTimes()
	err = times.Unmarshal(marshalledTimes)
	if err != nil {
		return nil
	}

	pathHash, err := pathHash(path)
	if err != nil {
		return nil
	}

	denyPath := filepath.Join(config.DenyDir(), pathHash)

	return &RC{path, allowPath, denyPath, times, config}
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
	if err = rc.times.Update(rc.allowPath); err != nil {
		return
	}
	if _, err = os.Stat(rc.denyPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	return os.Remove(rc.denyPath)
}

// Deny revokes the permission of the RC file to load
func (rc *RC) Deny() (err error) {
	if err = os.MkdirAll(filepath.Dir(rc.denyPath), 0755); err != nil {
		return
	}

	if err = os.WriteFile(rc.denyPath, []byte(rc.path+"\n"), 0644); /* #nosec G306 -- these deny files are not private */ err != nil {
		return
	}

	if _, err = os.Stat(rc.allowPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}

	return os.Remove(rc.allowPath)
}

type AllowStatus int

const (
	Allowed AllowStatus = iota
	NotAllowed
	Denied
)

// Allowed checks if the RC file has been granted loading
func (rc *RC) Allowed() AllowStatus {
	_, err := os.Stat(rc.denyPath)

	if err == nil {
		return Denied
	}

	// happy path is if this envrc has been explicitly allowed, O(1)ish common case
	_, err = os.Stat(rc.allowPath)

	if err == nil {
		return Allowed
	}

	// when whitelisting we want to be (path) absolutely sure we've not been duped with a symlink
	path, err := filepath.Abs(rc.path)
	// seems unlikely that we'd hit this, but have to handle it
	if err != nil {
		return NotAllowed
	}

	// exact whitelists are O(1)ish to check, so look there first
	if rc.config.WhitelistExact[path] {
		return Allowed
	}

	// finally we check if any of our whitelist prefixes match
	for _, prefix := range rc.config.WhitelistPrefix {
		if strings.HasPrefix(path, prefix) {
			return Allowed
		}
	}

	return NotAllowed
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
func (rc *RC) Load(previousEnv Env) (newEnv Env, err error) {
	config := rc.config
	wd := config.WorkDir
	direnv := config.SelfPath
	newEnv = previousEnv.Copy()
	newEnv[DIRENV_WATCHES] = rc.times.Marshal()
	defer func() {
		// Record directory changes even if load is disallowed or fails
		newEnv[DIRENV_DIR] = "-" + filepath.Dir(rc.path)
		newEnv[DIRENV_FILE] = rc.path
		newEnv[DIRENV_DIFF] = previousEnv.Diff(newEnv).Serialize()
	}()

	// Abort if the file is not allowed
	switch rc.Allowed() {
	case NotAllowed:
		err = fmt.Errorf(notAllowed, rc.Path())
		return
	case Allowed:
	case Denied:
		return
	}

	// Allow RC loads to be canceled with SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	// check what type of RC we're processing
	// use different exec method for each
	fn := "source_env"
	if filepath.Base(rc.path) == ".env" {
		fn = "dotenv"
	}

	// Set stdin based on the config
	var stdin *os.File
	if config.DisableStdin {
		stdin, err = os.Open(os.DevNull)
		if err != nil {
			return
		}
	} else {
		stdin = os.Stdin
	}

	prelude := ""
	if config.StrictEnv {
		prelude = "set -euo pipefail && "
	}

	// Non-Windows platforms will already use slashes. However, on Windows
	// backslashes are used by default which can result in unexpected escapes
	// like \b or \r in paths. Force slash usage to avoid issues on Windows.
	slashSeparatedPath := filepath.ToSlash(rc.Path())
	arg := fmt.Sprintf(
		`%seval "$("%s" stdlib)" && __main__ %s %s`,
		prelude,
		direnv,
		fn,
		BashEscape(slashSeparatedPath),
	)

	// G204: Subprocess launched with function call as argument or cmd arguments
	// #nosec
	cmd := exec.CommandContext(ctx, config.BashPath, "-c", arg)
	cmd.Dir = wd
	cmd.Env = newEnv.ToGoEnv()
	cmd.Stdin = stdin
	cmd.Stderr = os.Stderr

	var out []byte
	if out, err = cmd.Output(); err == nil && len(out) > 0 {
		var newEnv2 Env
		newEnv2, err = LoadEnvJSON(out)
		if err == nil {
			newEnv = newEnv2
		}
	}

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
	// Some broken filesystems like SSHFS return file information on stat() but
	// then cannot open the file. So we use os.Open.
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// Next, check that the file is a regular file.
	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
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

func pathHash(path string) (hash string, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte(path + "\n"))
	if err != nil {
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
	// G306: Expect WriteFile permissions to be 0600 or less
	// #nosec
	return os.WriteFile(allowPath, []byte(path+"\n"), 0644)
}

func findEnvUp(searchDir string, loadDotenv bool) (path string) {
	if loadDotenv {
		return findUp(searchDir, ".envrc", ".env")
	}
	return findUp(searchDir, ".envrc")
}

func findUp(searchDir string, fileNames ...string) (path string) {
	for _, dir := range eachDir(searchDir) {
		for _, fileName := range fileNames {
			path := filepath.Join(dir, fileName)
			if fileExists(path) {
				return path
			}
		}
	}
	return ""
}
