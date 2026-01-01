package cmd

import (
	"log"
	"os"
	"path"
	"strings"
)

// CmdPrune is `direnv prune`
var CmdPrune = &Cmd{
	Name:   "prune",
	Desc:   "Removes old allowed and required files",
	Action: actionWithConfig(cmdPruneAction),
}

func cmdPruneAction(_ Env, _ []string, config *Config) (err error) {
	var dir *os.File
	var fi os.FileInfo
	var dirList []string
	var envrc []byte

	// Track valid envrc paths for pruning required directory
	validEnvrcs := make(map[string]string) // pathHash -> envrcPath

	allowed := config.AllowDir()
	if dir, err = os.Open(allowed); err != nil {
		return err
	}
	defer func() {
		if err := dir.Close(); err != nil {
			log.Printf("Warning: failed to close directory: %v", err)
		}
	}()

	if dirList, err = dir.Readdirnames(0); err != nil {
		return err
	}

	for _, hash := range dirList {
		filename := path.Join(allowed, hash)
		if fi, err = os.Stat(filename); err != nil {
			return err
		}

		if !fi.IsDir() {
			if envrc, err = os.ReadFile(filename); err != nil {
				return err
			}
			envrcStr := strings.TrimSpace(string(envrc))

			// skip old files, w/o path inside
			if envrcStr == "" {
				continue
			}
			if !fileExists(envrcStr) {
				_ = os.Remove(filename)
			} else {
				// remove outdated hashes
				h, err := fileHash(envrcStr)
				if err != nil {
					return err
				}
				if h != hash {
					_ = os.Remove(filename)
				} else {
					// This envrc is still valid, track it
					if ph, err := pathHash(envrcStr); err == nil {
						validEnvrcs[ph] = envrcStr
					}
				}
			}

		}

	}

	// Prune orphaned and outdated allowed-required files
	return pruneAllowedRequiredDir(config, validEnvrcs)
}

func pruneAllowedRequiredDir(config *Config, validEnvrcs map[string]string) error {
	allowedRequiredDir := config.AllowedRequiredDir()
	dir, err := os.Open(allowedRequiredDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func() {
		if err := dir.Close(); err != nil {
			log.Printf("Warning: failed to close directory: %v", err)
		}
	}()

	dirList, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, envrcPathHash := range dirList {
		envrcPath, valid := validEnvrcs[envrcPathHash]
		if !valid {
			// Remove allowed-required directories that don't have a valid allowed envrc
			_ = os.RemoveAll(path.Join(allowedRequiredDir, envrcPathHash))
			continue
		}

		// Prune outdated allowed-required files within valid directories
		envrcDir := path.Dir(envrcPath)
		subdir := path.Join(allowedRequiredDir, envrcPathHash)
		if err := pruneAllowedRequiredFiles(subdir, envrcDir); err != nil {
			return err
		}
	}

	return nil
}

func pruneAllowedRequiredFiles(allowedRequiredSubdir, envrcDir string) error {
	dir, err := os.Open(allowedRequiredSubdir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func() {
		if err := dir.Close(); err != nil {
			log.Printf("Warning: failed to close directory: %v", err)
		}
	}()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, hash := range files {
		filename := path.Join(allowedRequiredSubdir, hash)
		content, err := os.ReadFile(filename)
		if err != nil {
			continue
		}
		relPath := strings.TrimSpace(string(content))

		absPath := path.Join(envrcDir, relPath)
		if !fileExists(absPath) {
			_ = os.Remove(filename)
		} else {
			// Check if hash is still valid
			h, err := fileHash(absPath)
			if err != nil || h != hash {
				_ = os.Remove(filename)
			}
		}
	}

	return nil
}
