package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
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

	var cmd *exec.Cmd

	if isPowerShellRC(rc.path, config) {
		if config.PwshPath == "" {
			err = fmt.Errorf("PowerShell file found but pwsh executable not available")
			return
		}

		// Safe execution strategy:
		// 1. Read windows.envrc content.
		// 2. Write a temporary wrapper script that dot-sources the file (supports functions) then emits JSON.
		// 3. Invoke pwsh with -File to avoid inline injection issues.
		// 4. Remove the temp wrapper after execution.

		envrcPath := rc.Path()
		contentBytes, readErr := os.ReadFile(envrcPath)
		if readErr != nil {
			err = readErr
			return
		}

		cacheDir := config.CacheDir
		if cacheDir == "" {
			err = fmt.Errorf("missing cache directory for PowerShell execution")
			return
		}

		// Use hash of path+mtime to avoid re-writing every time unnecessarily.
		stat, statErr := os.Stat(envrcPath)
		if statErr != nil {
			err = statErr
			return
		}
		contentHash := sha256.Sum256(contentBytes)
		hashInput := fmt.Sprintf("%s:%d:%x", envrcPath, stat.ModTime().UnixNano(), contentHash[:8])
		sha := sha256.Sum256([]byte(hashInput))
		wrapperName := fmt.Sprintf("direnv-pwsh-%x.ps1", sha[:8])
		wrapperPath := filepath.Join(cacheDir, wrapperName)

		// PowerShell wrapper now computes a delta (changed & removed) vs original env
		wrapperContent := fmt.Sprintf(`$ErrorActionPreference='Stop'
$ProgressPreference='SilentlyContinue'
try {
	$before = @{}
	Get-ChildItem env: | ForEach-Object { $before[$_.Name] = $_.Value }
	$envrcContent = @'
%s
'@
	# Execute the envrc content while suppressing non-error output streams so that the
	# only stdout produced by this wrapper is the final JSON payload. This prevents
	# Write-Host / Write-Output / warnings etc. from corrupting the JSON.
	# Streams: 1=Success, 3=Warning, 4=Verbose, 6=Information. Write-Host maps to Information.
	& { Invoke-Expression $envrcContent } 1>$null 3>$null 4>$null 6>$null
	$after = @{}
	Get-ChildItem env: | ForEach-Object { $after[$_.Name] = $_.Value }
	$changed = @{}
	foreach ($k in $after.Keys) { if (-not $before.ContainsKey($k) -or $before[$k] -ne $after[$k]) { $changed[$k] = $after[$k] } }
	$removed = @()
	foreach ($k in $before.Keys) { if (-not $after.ContainsKey($k)) { $removed += $k } }
	$result = @{ changed = $changed; removed = $removed }
	$json = $result | ConvertTo-Json -Compress
	[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
	Write-Output $json
} catch {
	$msg = "direnv pwsh error: $($_.Exception.Message)"
	[Console]::Error.WriteLine($msg)
	exit 1
}`, string(contentBytes))

		// Wrapper caching: reuse existing wrapper if contents match exact hash-generated name.
		// Since wrapperName already includes contentHash prefix via hashInput, if file exists we can skip rewrite.
		if _, statErr2 := os.Stat(wrapperPath); statErr2 != nil {
			if writeErr := os.WriteFile(wrapperPath, []byte(wrapperContent), 0600); writeErr != nil { // #nosec G304 - controlled path
				err = writeErr
				return
			}
		}

		cmd = exec.CommandContext(ctx, config.PwshPath, "-NoProfile", "-NonInteractive", "-File", wrapperPath)

		// Schedule cleanup of wrapper after execution attempt.
		defer func() {
			_ = os.Remove(wrapperPath)
		}()
	} else {
		// Execute bash .envrc or .env file (existing logic)
		fn := "source_env"
		if filepath.Base(rc.path) == ".env" {
			fn = "dotenv"
		}

		prelude := ""
		if config.StrictEnv {
			prelude = "set -euo pipefail && "
		}

		slashSeparatedPath := filepath.ToSlash(rc.Path())
		arg := fmt.Sprintf(
			`%seval "$("%s" stdlib)" && __main__ %s %s`,
			prelude,
			direnv,
			fn,
			BashEscape(slashSeparatedPath),
		)

		cmd = exec.CommandContext(ctx, config.BashPath, "-c", arg)
	}

	cmd.Dir = wd
	cmd.Env = newEnv.ToGoEnv()
	cmd.Stdin = stdin
	cmd.Stderr = os.Stderr

	var out []byte
	if out, err = cmd.Output(); err == nil && len(out) > 0 {
		if isPowerShellRC(rc.path, config) {
			// Parse delta JSON {"changed": {..}, "removed": [..]}
			type pwshDelta struct {
				Changed map[string]string `json:"changed"`
				Removed []string          `json:"removed"`
			}
			var delta pwshDelta
			if perr := json.Unmarshal(out, &delta); perr == nil {
				for k, v := range delta.Changed {
					newEnv[k] = v
				}
				for _, k := range delta.Removed {
					delete(newEnv, k)
				}
			} else {
				logError(config, fmt.Sprintf("pwsh delta parse failed: %v; attempting snapshot fallback", perr))
				// Fallback: try full snapshot parsing for backward-compatibility
				if newEnv2, jerr := LoadEnvJSON(out); jerr == nil {
					newEnv = newEnv2
				}
			}
		} else {
			var newEnv2 Env
			newEnv2, err = LoadEnvJSON(out)
			if err == nil {
				newEnv = newEnv2
			}
		}
	}

	return
}

// isPowerShellRC determines if an RC file should be executed as PowerShell
func isPowerShellRC(path string, config *Config) bool {
	if !config.EnablePwsh || config.PwshPath == "" {
		return false
	}
	return filepath.Base(path) == ".envrc.ps1"
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

func canonicalizeHashPath(path string) string {
	// On Windows, file paths are case-insensitive. Normalize to lower-case to
	// avoid generating distinct hashes for the same path differing only by case
	// (e.g. C:\Proj vs c:\proj). Also attempt to resolve symlinks so that
	// referencing a directory through a junction or symlink yields a stable
	// canonical path component for hashing.
	if runtime.GOOS == "windows" {
		// Lower-case drive letter and the rest
		path = strings.ToLower(path)
		if real, err := filepath.EvalSymlinks(path); err == nil {
			// EvalSymlinks returns path with native separators; lower-case again just in case
			path = strings.ToLower(real)
		}
	}
	return path
}

func fileHash(path string) (hash string, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}
	path = canonicalizeHashPath(path)

	fd, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() { _ = fd.Close() }()

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
	path = canonicalizeHashPath(path)

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
	if runtime.GOOS == "windows" {
		// Only support .envrc.ps1 (PowerShell) plus standard .envrc/.env
		path = findUp(searchDir, ".envrc.ps1")
		if path != "" {
			return path
		}
		if loadDotenv {
			return findUp(searchDir, ".envrc", ".env")
		}
		return findUp(searchDir, ".envrc")
	}
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
