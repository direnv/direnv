package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestWindowsEnvrcNoise ensures that extraneous Write-Host / Write-Output calls
// in a windows.envrc do not corrupt the JSON delta produced by the wrapper.
func TestWindowsEnvrcNoise(t *testing.T) {
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
	noisy := "Write-Host 'This should be suppressed'\r\nWrite-Output 'Also suppressed'\r\n$env:NOISE_VAR='ok'\r\n"
	if err := os.WriteFile(envrcPath, []byte(noisy), 0600); err != nil {
		t.Fatalf("write windows.envrc: %v", err)
	}
	rc, err := RCFromPath(envrcPath, config)
	if err != nil {
		t.Fatalf("RCFromPath: %v", err)
	}
	if err := rc.Allow(); err != nil {
		t.Fatalf("allow: %v", err)
	}
	base := GetEnv()
	newEnv, err := rc.Load(base)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if newEnv["NOISE_VAR"] != "ok" {
		t.Fatalf("expected NOISE_VAR=ok, got %q", newEnv["NOISE_VAR"])
	}
}
