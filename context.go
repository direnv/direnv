package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type LoadStatus uint

const (
	None LoadStatus = iota
	Loaded
	Unauthorized
	ScriptError
	InternalError
)

type Context struct {
	Env     Env
	Status  LoadStatus
	Error   error
	WorkDir string // Current directory
	ConfDir string
	ExecDir string
	RCDir   string
}

func LoadContext(env Env) (context *Context, err error) {
	context = &Context{
		Env:    env,
		Status: None,
	}

	context.ConfDir = env["DIRENV_CONFIG"]
	if context.ConfDir == "" {
		context.ConfDir = XdgConfigDir(env, "direnv")
	}
	if context.ConfDir == "" {
		err = fmt.Errorf("Couldn't find a configuration directory for direnv")
		return
	}

	context.ExecDir = env["DIRENV_LIBEXEC"]
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

func (self *Context) AllowDir() string {
	return filepath.Join(self.ConfDir, "allow")
}

func (self *Context) LoadedRC() *RC {
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

func (self *Context) FoundRC() *RC {
	return FindRC(self.WorkDir, self.AllowDir())
}

func (self *Context) EnvBackup() (Env, error) {
	if self.Env["DIRENV_BACKUP"] == "" {
		return nil, fmt.Errorf("DIRENV_BACKUP is empty")
	}
	return ParseEnv(self.Env["DIRENV_BACKUP"])
}
