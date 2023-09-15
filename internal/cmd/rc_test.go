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

func TestAllowFilePermission(t *testing.T) {
	rcPath := filepath.Join(t.TempDir(), ".envrc")
	allowPath := filepath.Join(t.TempDir(), "allow")
	_ = os.WriteFile(rcPath, nil, 0644) //nolint:gosec

	if err := allow(rcPath, allowPath); err != nil {
		t.Fatalf("allow failed: %v", err)
	}
	fi, err := os.Stat(rcPath)
	if err != nil {
		t.Fatal(err)
	}
	if perm := int(fi.Mode().Perm()); perm != 0600 {
		t.Fatalf("allow set incorrect permission: %#o", perm)
	}
}
