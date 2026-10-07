package cmd

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
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

// AllowStatus represents the permission status of an RC file.
type AllowStatus int

const (
	// Allowed indicates the RC file is permitted to load.
	Allowed AllowStatus = iota
	// NotAllowed indicates the RC file has not been granted permission.
	NotAllowed
	// Denied indicates the RC file has been explicitly denied.
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

	// A blocked or denied RC is never evaluated. If an allowed RC exists further
	// up, evaluate that one instead so the trusted environment above is kept.
	// The blocked RC is still recorded as DIRENV_FILE, and all the skipped RCs
	// are watched, so allowing or editing any of them triggers a reload.
	path := rc.path
	if status := rc.Allowed(); status != Allowed {
		allowed, skipped := rc.closestAllowed()

		watched := skipped
		if allowed != nil {
			watched = append(watched, allowed)
		}
		times := NewFileTimes()
		for _, s := range watched {
			for _, t := range *s.times.list {
				if err = times.NewTime(t.Path, t.Modtime, t.Exists); err != nil {
					return
				}
			}
		}
		newEnv[DIRENV_WATCHES] = times.Marshal()

		// When there is nothing to load, the returned error reports rc itself
		reported := skipped
		if allowed == nil {
			reported = skipped[1:]
		}
		for _, s := range reported {
			if s.Allowed() == NotAllowed {
				logError(config, "error "+notAllowed, s.Path())
			}
		}

		if allowed == nil {
			if status == NotAllowed {
				err = fmt.Errorf(notAllowed, rc.Path())
			}
			return
		}
		path = allowed.path
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
	if filepath.Base(path) == ".env" {
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
	slashSeparatedPath := filepath.ToSlash(path)
	arg := fmt.Sprintf(
		`%seval "$(%s stdlib)" && __main__ %s %s`,
		prelude,
		BashEscape(filepath.ToSlash(direnv)),
		fn,
		BashEscape(slashSeparatedPath),
	)

	// G204: Subprocess launched with function call as argument or cmd arguments
	// #nosec
	cmd := exec.CommandContext(ctx, config.BashPath, "-c", arg)
	cmd.Dir = wd
	cmd.Env = newEnv.ToGoEnv()
	cmd.Stdin = stdin
	// an *os.File is not copied through a pipe, so forked processes holding
	// it open do not block Wait
	cmd.Stderr = os.Stderr

	var stdout io.ReadCloser
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	var buf bytes.Buffer
	reader := bufio.NewReader(stdout)

	err = cmd.Start()
	if err != nil {
		return
	}

	// Stop at the end of the JSON output to not wait for forked
	// subprocesses forever
	for {
		line, readErr := reader.ReadBytes('\n')
		buf.Write(line)
		if readErr != nil || bytes.Equal(line, []byte("}\n")) {
			break
		}
	}
	_ = stdout.Close()

	if err = cmd.Wait(); err == nil && buf.Len() > 0 {
		var newEnv2 Env
		newEnv2, err = LoadEnvJSON(buf.Bytes())
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
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v", err)
		}
	}()

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
	if searchDir == "" {
		return ""
	}
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

// evaluated returns the RC that Load evaluates for rc: rc itself when it is
// allowed, otherwise the closest allowed RC above it. It also returns the RCs
// that get skipped.
func (rc *RC) evaluated() (allowed *RC, skipped []*RC) {
	if rc.Allowed() == Allowed {
		return rc, nil
	}
	return rc.closestAllowed()
}

// keepLoadedEnv handles moving from the loaded RC to the RC at path when both
// evaluate the same allowed RC, for example when entering a directory with a
// blocked .envrc. The environment is kept as it is, so the RC is not
// evaluated again; only DIRENV_FILE and the watches are updated and the newly
// skipped RCs are reported. It returns nil when the RC has to be loaded.
func keepLoadedEnv(loaded *RC, path string, env Env, config *Config) Env {
	if loaded.times.Check() != nil {
		return nil
	}
	next, err := RCFromPath(path, config)
	if err != nil {
		return nil
	}
	prevAllowed, prevSkipped := loaded.evaluated()
	nextAllowed, nextSkipped := next.evaluated()
	if prevAllowed == nil || nextAllowed == nil || prevAllowed.path != nextAllowed.path {
		return nil
	}

	// Stop watching the RCs that were skipped before, watch the new ones
	unwatched := make(map[string]bool)
	for _, s := range prevSkipped {
		for _, p := range []string{s.path, s.allowPath, s.denyPath} {
			if p, err = filepath.Abs(p); err == nil {
				unwatched[filepath.Clean(p)] = true
			}
		}
	}
	current := NewFileTimes()
	if err = current.Unmarshal(env[DIRENV_WATCHES]); err != nil {
		return nil
	}
	times := NewFileTimes()
	for _, t := range *current.list {
		if !unwatched[t.Path] {
			if err = times.NewTime(t.Path, t.Modtime, t.Exists); err != nil {
				return nil
			}
		}
	}
	for _, s := range nextSkipped {
		for _, t := range *s.times.list {
			if err = times.NewTime(t.Path, t.Modtime, t.Exists); err != nil {
				return nil
			}
		}
		if s.Allowed() == NotAllowed {
			logError(config, "error "+notAllowed, s.Path())
		}
	}

	newEnv := env.Copy()
	newEnv[DIRENV_WATCHES] = times.Marshal()
	newEnv[DIRENV_DIR] = "-" + filepath.Dir(next.path)
	newEnv[DIRENV_FILE] = next.path
	return newEnv
}

// closestAllowed looks for the closest RC that is allowed to load, starting
// next to rc and going up. It also returns the RCs that were skipped on the
// way, rc included.
func (rc *RC) closestAllowed() (allowed *RC, skipped []*RC) {
	skipped = []*RC{rc}
	rcPath, err := filepath.Abs(rc.path)
	if err != nil {
		return nil, skipped
	}

	fileNames := []string{".envrc"}
	if rc.config.LoadDotenv {
		fileNames = append(fileNames, ".env")
	}

	for _, dir := range eachDir(filepath.Dir(rcPath)) {
		for _, fileName := range fileNames {
			path := filepath.Join(dir, fileName)
			if path == rcPath || !fileExists(path) {
				continue
			}
			candidate, err := RCFromPath(path, rc.config)
			if err != nil {
				continue
			}
			if candidate.Allowed() == Allowed {
				return candidate, skipped
			}
			skipped = append(skipped, candidate)
		}
	}
	return nil, skipped
}
