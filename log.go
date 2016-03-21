package main

import (
	"fmt"
	"log"
)

const (
	debugLogFormat   = "DBG-direnv: %s"
	defaultLogFormat = "direnv: %s"
	errorLogFormat   = "\033[31mdirenv: %s\033[0m"
)

var debugging bool

func setupLogging(env Env) {
	log.SetFlags(0)
	log.SetPrefix("")
	if val, ok := env[DIRENV_DEBUG]; ok == true && val == "1" {
		debugging = true
		log.SetFlags(log.Ltime)
		log.SetPrefix("direnv: ")
	}
}

func log_error(msg string, a ...interface{}) {
	logMsg(errorLogFormat, msg, a...)
}

func log_status(env Env, msg string, a ...interface{}) {
	format, ok := env["DIRENV_LOG_FORMAT"]
	if !ok {
		format = defaultLogFormat
	}
	if format != "" {
		logMsg(format, msg, a...)
	}
}

func log_debug(msg string, a ...interface{}) {
	if !debugging {
		return
	}
	defer log.SetFlags(log.Flags())
	log.SetFlags(log.Flags() | log.Lshortfile)
	msg = fmt.Sprintf(msg, a...)
	log.Output(2, msg)
}

func logMsg(format, msg string, a ...interface{}) {
	defer log.SetFlags(log.Flags())
	defer log.SetPrefix(log.Prefix())
	log.SetFlags(0)
	log.SetPrefix("")

	msg = fmt.Sprintf(format+"\n", msg)
	log.Printf(msg, a...)
}
