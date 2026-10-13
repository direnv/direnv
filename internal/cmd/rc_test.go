package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSomething(t *testing.T) {
	paths := eachDir("/foo/b//bar/")
	if len(paths) != 4 {
		t.Fail()
	}
	// TODO: fix me for windows
	if runtime.GOOS != "windows" {
		paths = eachDir("/")
		if len(paths) != 1 && paths[0] != "/" {
			t.Fail()
		}
	}
}

func TestClosestAllowed(t *testing.T) {
	root := t.TempDir()
	config := &Config{DataDir: filepath.Join(root, "data")}

	// root/.envrc is allowed, a/.envrc is blocked, a/b/.envrc is denied
	newRC := func(dir string) *RC {
		t.Helper()
		dir = filepath.Join(root, dir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, ".envrc")
		if err := os.WriteFile(path, []byte("export FOO=1\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		rc, err := RCFromPath(path, config)
		if err != nil {
			t.Fatal(err)
		}
		return rc
	}
	top := newRC(".")
	middle := newRC("a")
	nearest := newRC("a/b")
	if err := top.Allow(); err != nil {
		t.Fatal(err)
	}
	if err := nearest.Deny(); err != nil {
		t.Fatal(err)
	}

	allowed, skipped := nearest.closestAllowed()
	if allowed == nil || allowed.Path() != top.Path() {
		t.Fatalf("expected %s to be loaded, got %v", top.Path(), allowed)
	}
	if len(skipped) != 2 || skipped[0].Path() != nearest.Path() || skipped[1].Path() != middle.Path() {
		t.Fatalf("expected %s and %s to be skipped, got %v", nearest.Path(), middle.Path(), skipped)
	}

	// Without any allowed RC above, there is nothing to load
	if err := top.Deny(); err != nil {
		t.Fatal(err)
	}
	if allowed, _ := nearest.closestAllowed(); allowed != nil {
		t.Fatalf("expected nothing to be loaded, got %s", allowed.Path())
	}
}
