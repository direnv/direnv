package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		if bashPath != "" {
			config.BashPath = bashPath
		} else if config.BashPath, err = exec.LookPath("bash"); err != nil {
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
		log_debug("RCDir is blank - loadedRC is nil")
		return nil
	}
	rcPath := filepath.Join(self.RCDir, ".envrc")

	times_string := self.Env[DIRENV_WATCHES]

	return RCFromEnv(rcPath, times_string)
}

func (self *Config) FindRC() *RC {
	return FindRC(self.WorkDir, self.AllowDir())
}

func (self *Config) EnvDiff() (*EnvDiff, error) {
	if self.Env[DIRENV_DIFF] == "" {
		if self.Env[DIRENV_WATCHES] == "" {
			return self.Env.Diff(self.Env), nil
		} else {
			return nil, fmt.Errorf("DIRENV_DIFF is empty")
		}
	}
	return LoadEnvDiff(self.Env[DIRENV_DIFF])
}
