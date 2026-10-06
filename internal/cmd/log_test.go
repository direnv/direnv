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
