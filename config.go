package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Env      Env
	WorkDir  string // Current directory
	ConfDir  string
	SelfPath string
	BashPath string
	RCDir    string
}

func LoadConfig(env Env) (config *Config, err error) {
	config = &Config{
		Env: env,
	}

	config.ConfDir = env[DIRENV_CONFIG]
	if config.ConfDir == "" {
		config.ConfDir = XdgConfigDir(env, "direnv")
	}
	if config.ConfDir == "" {
		err = fmt.Errorf("Couldn't find a configuration directory for direnv")
		return
	}

	var exePath string
	if exePath, err = exec.LookPath(os.Args[0]); err != nil {
		err = fmt.Errorf("LoadConfig() Lookpath failed: %q", err)
		return
	}
	if exePath, err = filepath.EvalSymlinks(exePath); err != nil {
		err = fmt.Errorf("LoadConfig() symlink resolution: %q", err)
		return
	}
	exePath = strings.Replace(exePath, "\\", "/", -1)
	config.SelfPath = exePath

	config.BashPath = env[DIRENV_BASH]
	if config.BashPath == "" {
		if config.BashPath, err = exec.LookPath("bash"); err != nil {
			err = fmt.Errorf("Can't find bash: %q", err)
			return
		}
	}

	if config.WorkDir, err = os.Getwd(); err != nil {
		err = fmt.Errorf("LoadConfig() Getwd failed: %q", err)
		return
	}

	config.RCDir = env[DIRENV_DIR]
	if len(config.RCDir) > 0 && config.RCDir[0:1] == "-" {
		config.RCDir = config.RCDir[1:]
	}

	return
}

func (self *Config) AllowDir() string {
	return filepath.Join(self.ConfDir, "allow")
}

func (self *Config) LoadedRC() *RC {
	if self.RCDir == "" {
		return nil
	}
	rcPath := filepath.Join(self.RCDir, ".envrc")

	mtime, err := strconv.ParseInt(self.Env[DIRENV_MTIME], 10, 64)
	if err != nil {
		return nil
	}

	return RCFromEnv(rcPath, mtime)
}

func (self *Config) FindRC() *RC {
	return FindRC(self.WorkDir, self.AllowDir())
}

func (self *Config) EnvDiff() (*EnvDiff, error) {
	if self.Env[DIRENV_DIFF] == "" {
		return nil, fmt.Errorf("DIRENV_DIFF is empty")
	}
	return LoadEnvDiff(self.Env[DIRENV_DIFF])
}
