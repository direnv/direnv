package main

import (
	"os"
	"strings"
)

type emacs int

var EMACS emacs

func (x emacs) Hook() string {
	log_error("this feature is not supported. Install the direnv.el plugin instead.")
	os.Exit(1)
	return ""
}

func (x emacs) Escape(str string) string {
	str = strings.Replace(str, "\\", "\\\\", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	return "\"" + str + "\""
}

func (x emacs) Export(key, value string) string {
	return "(setenv " + x.Escape(key) + " " + x.Escape(value) + ")\n"
}

func (x emacs) Unset(key string) string {
	return "(setenv " + x.Escape(key) + " nil)\n"
}
