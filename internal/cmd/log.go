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

func logDebug(msg string, a ...interface{}) {
	if !debugging {
		return
	}
	defer log.SetFlags(log.Flags())
	log.SetFlags(log.Flags() | log.Lshortfile)
	msg = fmt.Sprintf(msg, a...)
	_ = log.Output(2, msg)
}

func logMsg(format, msg string, a ...interface{}) {
	defer log.SetFlags(log.Flags())
	defer log.SetPrefix(log.Prefix())
	log.SetFlags(0)
	log.SetPrefix("")

	msg = fmt.Sprintf(format+"\n", msg)
	log.Printf(msg, a...)
}
