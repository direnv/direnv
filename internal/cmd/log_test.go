package cmd

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// captureLog captures log output during fn execution
func captureLog(fn func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)
	fn()
	return buf.String()
}

// Regression test for https://github.com/direnv/direnv/issues/1551
// When LogColor is false (dumb terminal), no ANSI escape sequences should be output.
func TestLogErrorNoColorOnDumbTerminal(t *testing.T) {
	c := &Config{LogColor: false}
	output := captureLog(func() {
		logError(c, "test error %s", "msg")
	})
	if strings.Contains(output, "\033[") {
		t.Errorf("logError with LogColor=false should not contain ANSI escapes, got: %q", output)
	}
	if !strings.Contains(output, "test error msg") {
		t.Errorf("logError should contain the message, got: %q", output)
	}
}

func TestLogErrorWithColor(t *testing.T) {
	c := &Config{LogColor: true}
	output := captureLog(func() {
		logError(c, "test error %s", "msg")
	})
	if !strings.Contains(output, "\033[31m") {
		t.Errorf("logError with LogColor=true should contain red ANSI escape, got: %q", output)
	}
}

func TestLogStatusNoColorOnDumbTerminal(t *testing.T) {
	c := &Config{LogColor: false, LogFormat: "direnv: %s"}
	output := captureLog(func() {
		logStatus(c, "loading .envrc")
	})
	if strings.Contains(output, "\033[") {
		t.Errorf("logStatus with LogColor=false should not contain ANSI escapes, got: %q", output)
	}
	if !strings.Contains(output, "loading .envrc") {
		t.Errorf("logStatus should contain the message, got: %q", output)
	}
}

func TestLogStatusWithColor(t *testing.T) {
	c := &Config{LogColor: true, LogFormat: "direnv: %s"}
	output := captureLog(func() {
		logStatus(c, "loading .envrc")
	})
	if !strings.Contains(output, "\033[0m") {
		t.Errorf("logStatus with LogColor=true should contain reset ANSI escape, got: %q", output)
	}
}
