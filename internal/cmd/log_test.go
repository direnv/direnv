package cmd

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Writer()
	prevFlags := log.Flags()
	prevPrefix := log.Prefix()
	log.SetOutput(&buf)
	log.SetFlags(0)
	log.SetPrefix("")
	defer func() {
		log.SetOutput(prev)
		log.SetFlags(prevFlags)
		log.SetPrefix(prevPrefix)
	}()
	fn()
	return buf.String()
}

func TestLogErrorColor(t *testing.T) {
	c := &Config{LogFormat: defaultLogFormat, LogColor: true}
	out := captureLog(t, func() { logError(c, "oops") })
	if !strings.Contains(out, errorColor) {
		t.Fatalf("expected red escape when LogColor=true, got %q", out)
	}
	if !strings.Contains(out, "direnv: oops") {
		t.Fatalf("expected message, got %q", out)
	}
	if !strings.Contains(out, clearColor) {
		t.Fatalf("expected clear escape when LogColor=true, got %q", out)
	}
}

func TestLogErrorNoColor(t *testing.T) {
	c := &Config{LogFormat: defaultLogFormat, LogColor: false}
	out := captureLog(t, func() { logError(c, "oops") })
	if strings.Contains(out, "\033[") {
		t.Fatalf("expected no escape codes when LogColor=false, got %q", out)
	}
	if strings.TrimSpace(out) != "direnv: oops" {
		t.Fatalf("unexpected output %q", out)
	}
}

func TestLogStatusColor(t *testing.T) {
	c := &Config{LogFormat: defaultLogFormat, LogColor: true}
	out := captureLog(t, func() { logStatus(c, "hello") })
	if !strings.HasPrefix(out, clearColor) {
		t.Fatalf("expected reset prefix when LogColor=true, got %q", out)
	}
	if !strings.Contains(out, "direnv: hello") {
		t.Fatalf("expected message, got %q", out)
	}
}

func TestLogStatusNoColor(t *testing.T) {
	c := &Config{LogFormat: defaultLogFormat, LogColor: false}
	out := captureLog(t, func() { logStatus(c, "hello") })
	if strings.Contains(out, "\033[") {
		t.Fatalf("expected no escape codes when LogColor=false, got %q", out)
	}
	if strings.TrimSpace(out) != "direnv: hello" {
		t.Fatalf("unexpected output %q", out)
	}
}

func TestLoadConfigLogColorRespectsTERM(t *testing.T) {
	tmp := t.TempDir()
	confDir := tmp + "/direnv"
	cacheDir := tmp + "/cache"
	dataDir := tmp + "/data"
	for _, d := range []string{confDir, cacheDir, dataDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	prevTERM := os.Getenv("TERM")
	t.Cleanup(func() { _ = os.Setenv("TERM", prevTERM) })

	env := Env{
		"DIRENV_CONFIG":  confDir,
		"XDG_CACHE_HOME": cacheDir,
		"XDG_DATA_HOME":  dataDir,
	}

	_ = os.Setenv("TERM", "xterm-256color")
	cfg, err := LoadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.LogColor {
		t.Fatalf("expected LogColor=true without toml when TERM is not dumb")
	}

	_ = os.Setenv("TERM", "dumb")
	cfg, err = LoadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.LogColor {
		t.Fatalf("expected LogColor=false when TERM=dumb")
	}

	if err := os.WriteFile(confDir+"/direnv.toml", []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = os.Setenv("TERM", "xterm-256color")
	cfg, err = LoadConfig(env)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.LogColor {
		t.Fatalf("expected LogColor=true with toml when TERM is not dumb")
	}
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	prev := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = prev }()
	fn()
	_ = w.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

// The error returned by a command is printed by Main, not by logError, so it
// has to follow the same TERM rule (issue #1509).
func TestMainErrorColorRespectsTERM(t *testing.T) {
	// Main unsets PWD; restore it when the test is done.
	t.Setenv("PWD", os.Getenv("PWD"))

	runMain := func() string {
		return captureStderr(t, func() {
			if err := Main(Env{}, []string{"direnv", "no-such-command"}, "", "", ""); err == nil {
				t.Fatal("expected an error for an unknown command")
			}
		})
	}

	t.Setenv("TERM", "dumb")
	out := runMain()
	if strings.Contains(out, "\033[") {
		t.Fatalf("expected no escape codes when TERM=dumb, got %q", out)
	}
	if strings.TrimSpace(out) != `direnv: error command "direnv no-such-command" not found` {
		t.Fatalf("unexpected output %q", out)
	}

	t.Setenv("TERM", "xterm-256color")
	out = runMain()
	if !strings.HasPrefix(out, errorColor) || !strings.Contains(out, clearColor) {
		t.Fatalf("expected colored error when TERM is not dumb, got %q", out)
	}
}
