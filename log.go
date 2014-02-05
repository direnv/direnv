package main

import (
	"fmt"
	"os"
)

const (
	defaultLogFormat = "direnv: %s"
	errorLogFormat   = "\033[31mdirenv: %s\033[0m"
)

var logFormat = ""

func formatLog(msg string) string {
	if logFormat == "" {
		logFormat = os.Getenv("DIRENV_LOG_FORMAT")
		if logFormat == "" {
			logFormat = defaultLogFormat
		}
	}
	return fmt.Sprintf(logFormat, msg)
}

func log_error(msg string, a ...interface{}) {
	msg = fmt.Sprintf(errorLogFormat, msg)
	log(msg, a...)
}

func log_status(msg string, a ...interface{}) {
	log(formatLog(msg), a...)
}

func log(msg string, a ...interface{}) {
	msg = fmt.Sprintf(msg, a...)
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}
