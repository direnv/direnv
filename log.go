package main

import (
	"fmt"
	"os"
)

const (
	GreyCode  = "\033[30m"
	RedCode   = "\033[31m"
	GreenCode = "\033[32m"
	CyanCode  = "\033[36m"
	PlainCode = "\033[0m"
	ResetCode = "\033[0m"
)

func log(msg string, a ...interface{}) {
	log_color(PlainCode, msg, a...)
}

func log_error(msg string, a ...interface{}) {
	log_color(RedCode, msg, a...)
}

func log_status(msg string, a ...interface{}) {
	log_color(GreyCode, msg, a...)
}

func log_color(color string, msg string, a ...interface{}) {
	msg = fmt.Sprintf(msg, a...)
	fmt.Fprintf(os.Stderr, "%sdirenv: %s%s\n", color, msg, ResetCode)
}
