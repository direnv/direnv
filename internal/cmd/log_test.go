package cmd

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func captureLogOutput(t *testing.T, fn func()) string {
	t.Helper()

	var buf bytes.Buffer
	origWriter := log.Writer()
	origFlags := log.Flags()
	origPrefix := log.Prefix()

	log.SetOutput(&buf)
	log.SetFlags(0)
	log.SetPrefix("")
	t.Cleanup(func() {
		log.SetOutput(origWriter)
		log.SetFlags(origFlags)
		log.SetPrefix(origPrefix)
	})

	fn()

	return buf.String()
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	origStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = origStderr
	})

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("closing stderr writer failed: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("reading captured stderr failed: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("closing stderr reader failed: %v", err)
	}

	return buf.String()
}

func TestLoadConfigSetsLogColorFromTERM(t *testing.T) {
	t.Run("ansi terminal", func(t *testing.T) {
		baseDir := t.TempDir()
		cfg, err := LoadConfig(Env{
			"DIRENV_CONFIG": baseDir + "/config",
			"HOME":          baseDir,
			"TERM":          "xterm-256color",
		})
		if err != nil {
			t.Fatalf("LoadConfig returned error: %v", err)
		}
		if !cfg.LogColor {
			t.Fatalf("expected LogColor to be enabled for ANSI-capable terminals")
		}
	})

	t.Run("dumb terminal", func(t *testing.T) {
		baseDir := t.TempDir()
		cfg, err := LoadConfig(Env{
			"DIRENV_CONFIG": baseDir + "/config",
			"HOME":          baseDir,
			"TERM":          "dumb",
		})
		if err != nil {
			t.Fatalf("LoadConfig returned error: %v", err)
		}
		if cfg.LogColor {
			t.Fatalf("expected LogColor to be disabled for dumb terminals")
		}
	})
}

func TestMainSkipsAnsiForTopLevelErrorsOnDumbTerminals(t *testing.T) {
	out := captureStderr(t, func() {
		err := Main(Env{"TERM": "dumb"}, []string{"direnv", "export", "nosuchshell"}, "", "", "")
		if err == nil {
			t.Fatal("expected Main to return an error for an unknown shell")
		}
	})
	if strings.Contains(out, "\033[") {
		t.Fatalf("expected top-level error without ANSI escapes, got %q", out)
	}
}

func TestLogStatusUsesAnsiOnlyWhenEnabled(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		out := captureLogOutput(t, func() {
			logStatus(&Config{LogColor: true, LogFormat: defaultLogFormat}, "loading test")
		})
		if !strings.Contains(out, clearColor) {
			t.Fatalf("expected status log to reset color when enabled, got %q", out)
		}
	})

	t.Run("disabled", func(t *testing.T) {
		out := captureLogOutput(t, func() {
			logStatus(&Config{LogColor: false, LogFormat: defaultLogFormat}, "loading test")
		})
		if strings.Contains(out, "\033[") {
			t.Fatalf("expected status log without ANSI escapes, got %q", out)
		}
	})
}

func TestLogErrorUsesAnsiOnlyWhenEnabled(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		out := captureLogOutput(t, func() {
			logError(&Config{LogColor: true}, "boom")
		})
		if !strings.Contains(out, errorColor) {
			t.Fatalf("expected error log to include ANSI color when enabled, got %q", out)
		}
	})

	t.Run("disabled", func(t *testing.T) {
		out := captureLogOutput(t, func() {
			logError(&Config{LogColor: false}, "boom")
		})
		if strings.Contains(out, "\033[") {
			t.Fatalf("expected error log without ANSI escapes, got %q", out)
		}
	})
}
