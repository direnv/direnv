package main

import (
	"os"
	"os/exec"
	"path/filepath"
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
	ExecDir string
	RC      *RC
}

var CONTEXT *Context

func LoadContext(env Env) (context *Context) {
	context = &Context{
		Env:    env,
		Status: None,
	}

	fail := func(err error) *Context {
		context.Error = err
		context.Status = InternalError
		return context
	}

	context.ExecDir = env["DIRENV_LIBEXEC"]
	if context.ExecDir == "" {
		exePath, err := exec.LookPath(os.Args[0])
		if err != nil {
			return fail(err)
		}

		exePath, err = filepath.EvalSymlinks(exePath)
		if err != nil {
			return fail(err)
		}

		context.ExecDir = filepath.Dir(exePath)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fail(err)
	}

	context.RC = FindRC(wd)

	return
}

