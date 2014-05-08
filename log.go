package main

import (
	"fmt"
	"os"
)

const (
	defaultLogFormat = "direnv: %s"
	errorLogFormat   = "\033[31mdirenv: %s\033[0m"
)

func log_error(msg string, a ...interface{}) {
	log(errorLogFormat, msg, a...)
}

func log_status(env Env, msg string, a ...interface{}) {
	format, ok := env["DIRENV_LOG_FORMAT"]
	if !ok {
		format = defaultLogFormat
	}
	if format != "" {
		log(format, msg, a...)
	}
}

func log(format, msg string, a ...interface{}) {
	msg = fmt.Sprintf(format+"\n", msg)
	fmt.Fprintf(os.Stderr, msg, a...)
}
