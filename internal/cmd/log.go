package cmd

import (
	"fmt"
	"log"
	"strings"
)

const (
	defaultLogFormat = "direnv: %s"
	errorColor       = "\033[31m"
	clearColor       = "\033[0m"
)

var debugging bool

func setupLogging(env Env) {
	log.SetFlags(0)
	log.SetPrefix("")
	if val, ok := env[DIRENV_DEBUG]; ok && (val == "1" || strings.EqualFold(val, "true")) {
		debugging = true
		log.SetFlags(log.Ltime)
		log.SetPrefix("direnv: ")
	}
}

func logError(c *Config, msg string, a ...interface{}) {
	if c.LogColor {
		logMsg(defaultLogFormat, msg, a...)
	} else {
		logMsg(errorColor+defaultLogFormat+clearColor, msg, a...)
	}
}

func logStatus(c *Config, msg string, a ...interface{}) {
	format := c.LogFormat
	shouldLog := true
	if c.LogFilter != nil {
		shouldLog = c.LogFilter.MatchString(msg)
	}
	if shouldLog && format != "" {
		if c.LogColor {
			logMsg(format, msg, a...)
		} else {
			logMsg(fmt.Sprintf("%s%s", clearColor, format), msg, a...)
		}
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
