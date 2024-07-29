package cmd

import (
	"fmt"
	"log"
	"os"
)

const (
	defaultLogFormat        = "direnv: %s"
	errorLogFormat          = defaultLogFormat
	errorLogFormatWithColor = "\033[31mdirenv: %s\033[0m"
)

var debugging bool
var noColor = os.Getenv("TERM") == "dumb"

func setupLogging(env Env) {
	log.SetFlags(0)
	log.SetPrefix("")
	if val, ok := env[DIRENV_DEBUG]; ok && val == "1" {
		debugging = true
		log.SetFlags(log.Ltime)
		log.SetPrefix("direnv: ")
	}
}

func logError(msg string, a ...interface{}) {
	if noColor {
		logMsg(errorLogFormat, msg, a...)
	} else {
		logMsg(errorLogFormatWithColor, msg, a...)
	}
}

func logStatus(env Env, msg string, a ...interface{}) {
	format, ok := env["DIRENV_LOG_FORMAT"]
	if !ok {
		format = defaultLogFormat
	}
	if format != "" {
		logMsg(format, msg, a...)
	}
}

func logDebugFormat(msg string, a ...interface{}) {
	defer log.SetFlags(log.Flags())
	log.SetFlags(log.Flags() | log.Lshortfile)
	msg = fmt.Sprintf(msg, a...)
	_ = log.Output(2, msg)
}

func logDebug(msg string, a ...interface{}) {
	if !debugging {
		return
	}
	logDebugFormat(msg, a...)
}

// logDebugInit is an alternative to logDebug that can be called before
// setupLogging is called.
// NOTE: logDebug depends on variables that are set by setupLogging, which in
// turn depends on the existence of an Env object. Sometimes we need debug
// logging before we have an Env object (see the cygpath implementation).
func logDebugInit(msg string, a ...interface{}) {
	// aka !debugging (see setupLogging and logDebug)
	if debug := os.Getenv("DIRENV_DEBUG"); debug != "1" {
		return
	}
	logDebugFormat(msg, a...)
}

func logMsg(format, msg string, a ...interface{}) {
	defer log.SetFlags(log.Flags())
	defer log.SetPrefix(log.Prefix())
	log.SetFlags(0)
	log.SetPrefix("")

	msg = fmt.Sprintf(format+"\n", msg)
	log.Printf(msg, a...)
}
