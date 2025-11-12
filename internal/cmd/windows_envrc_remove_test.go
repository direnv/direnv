package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestWindowsEnvrcRemoval verifies that removing a variable between loads unsets it via delta 'removed'.
func TestWindowsEnvrcRemoval(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows only")
	}
	env := GetEnv()
	config, err := LoadConfig(env)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if !config.EnablePwsh || config.PwshPath == "" {
		t.Skip("pwsh not available or disabled")
	}

	tmpDir := t.TempDir()
	envrcPath := filepath.Join(tmpDir, ".envrc.ps1")

	write := func(content string) {
		if err := os.WriteFile(envrcPath, []byte(content), 0600); err != nil {
			t.Fatalf("write windows.envrc: %v", err)
		}
	}

	// Initial version sets VAR1 and VAR2
	write("$env:VAR1='one'\r\n$env:VAR2='two'\r\n")
	rc, err := RCFromPath(envrcPath, config)
	if err != nil {
		t.Fatalf("RCFromPath: %v", err)
	}
	if err := rc.Allow(); err != nil {
		t.Fatalf("allow: %v", err)
	}
	base := GetEnv()
	after1, err := rc.Load(base)
	if err != nil {
		t.Fatalf("first load: %v", err)
	}
	if after1["VAR1"] != "one" || after1["VAR2"] != "two" {
		t.Fatalf("vars not set: %+v", after1)
	}

	// Second version explicitly removes VAR2 and changes VAR1
	write("$env:VAR1='one-mod'\r\n$env:VAR2=$null\r\nRemove-Item Env:VAR2 -ErrorAction SilentlyContinue\r\n")
	// Ensure mtime difference
	time.Sleep(20 * time.Millisecond)
	rc2, err := RCFromPath(envrcPath, config)
	if err != nil {
		t.Fatalf("RCFromPath second: %v", err)
	}
	// Already allowed; Allow again is harmless
	if err := rc2.Allow(); err != nil {
		t.Fatalf("allow second: %v", err)
	}
	after2, err := rc2.Load(after1)
	if err != nil {
		t.Fatalf("second load: %v", err)
	}
	if after2["VAR1"] != "one-mod" {
		t.Errorf("VAR1 expected one-mod got %q", after2["VAR1"])
	}
	if _, ok := after2["VAR2"]; ok {
		t.Errorf("VAR2 should be removed")
	}
}
