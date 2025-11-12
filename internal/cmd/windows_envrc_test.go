package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestWindowsEnvrc ensures that a windows.envrc file is detected and executed via PowerShell
// Only runs on Windows when pwsh is available.
func TestWindowsEnvrc(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows.envrc only applies on Windows")
	}

	// Ensure pwsh is discoverable in config
	env := GetEnv()
	config, err := LoadConfig(env)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if config.PwshPath == "" {
		t.Skip("pwsh not found in PATH; skipping PowerShell test")
	}

	// Create a temp directory to host the PowerShell envrc (.envrc.ps1 preferred name)
	tmpDir := t.TempDir()
	// Use separate lines without embedded \n to avoid Invoke-Expression parse errors
	rcContent := "$env:TEST_HELLO = \"pwsh-world\"\r\n$env:TEST_NUMBER = \"42\"\r\n"
	envrcPath := filepath.Join(tmpDir, ".envrc.ps1")
	if err := os.WriteFile(envrcPath, []byte(rcContent), 0600); err != nil {
		t.Fatalf("failed to write windows.envrc: %v", err)
	}

	// Load RC
	// Read the file to ensure no lingering writer handle
	if _, rerr := os.ReadFile(envrcPath); rerr != nil {
		t.Fatalf("read windows.envrc failed: %v", rerr)
	}

	rc, err := RCFromPath(envrcPath, config)
	if err != nil {
		t.Fatalf("RCFromPath error: %v", err)
	}

	// Mark as allowed (simulate direnv allow)
	if err := rc.Allow(); err != nil {
		t.Fatalf("rc.Allow error: %v", err)
	}

	baseEnv := GetEnv()
	newEnv, err := rc.Load(baseEnv)
	if err != nil {
		t.Fatalf("rc.Load error: %v", err)
	}

	if newEnv["TEST_HELLO"] != "pwsh-world" {
		t.Errorf("TEST_HELLO not set, got: %q", newEnv["TEST_HELLO"])
	}
	if newEnv["TEST_NUMBER"] != "42" {
		t.Errorf("TEST_NUMBER not set, got: %q", newEnv["TEST_NUMBER"])
	}

	if err := attemptRemove(envrcPath); err != nil {
		t.Logf("warning: attemptRemove windows.envrc: %v", err)
	}
}

// attemptRemove tries to delete a file with retries to mitigate transient locks on Windows
func attemptRemove(path string) error {
	var lastErr error
	for i := 0; i < 20; i++ {
		if err := os.Remove(path); err == nil || os.IsNotExist(err) {
			return nil
		} else {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
		}
	}
	return lastErr
}
