package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Config struct {
	Env     Env
	WorkDir string // Current directory
	ConfDir string
	ExecDir string
	RCDir   string
}

func LoadConfig(env Env) (context *Config, err error) {
	context = &Config{
		Env: env,
	}

	context.ConfDir = env["DIRENV_CONFIG"]
	if context.ConfDir == "" {
		context.ConfDir = XdgConfigDir(env, "direnv")
	}
	if context.ConfDir == "" {
		err = fmt.Errorf("Couldn't find a configuration directory for direnv")
		return
	}

	//context.ExecDir = env["DIRENV_LIBEXEC"]
	if context.ExecDir == "" {
		var exePath string
		if exePath, err = exec.LookPath(os.Args[0]); err != nil {
			return
		}

		if exePath, err = filepath.EvalSymlinks(exePath); err != nil {
			return
		}

		context.ExecDir = filepath.Dir(exePath)
	}

	if context.WorkDir, err = os.Getwd(); err != nil {
		return
	}

	context.RCDir = env["DIRENV_DIR"]
	if len(context.RCDir) > 0 && context.RCDir[0:1] == "-" {
		context.RCDir = context.RCDir[1:]
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

	mtime, err := strconv.ParseInt(self.Env["DIRENV_MTIME"], 10, 64)
	if err != nil {
		return nil
	}

	if self.Env["DIRENV_HASH"] == "" {
		return nil
	}
	hash := self.Env["DIRENV_HASH"]

	return RCFromEnv(rcPath, mtime, hash, self.AllowDir())
}

func (self *Config) FoundRC() *RC {
	return FindRC(self.WorkDir, self.AllowDir())
}

func (self *Config) EnvBackup() (Env, error) {
	if self.Env["DIRENV_BACKUP"] == "" {
		return nil, fmt.Errorf("DIRENV_BACKUP is empty")
	}
	return ParseEnv(self.Env["DIRENV_BACKUP"])
}
