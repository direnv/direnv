//go:build windows
// +build windows

package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// Test that fileHash and pathHash are stable across drive-letter / case differences
func TestWindowsPathCaseHashStability(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	// Create a temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, ".envrc.ps1")
	if err := os.WriteFile(tmpFile, []byte("Write-Host 'hi'"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	// Build two path variants that differ only in case. On Windows filepath.Abs keeps drive letter casing from input.
	absPath, err := filepath.Abs(tmpFile)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if len(absPath) < 2 || absPath[1] != ':' {
		t.Skipf("unexpected path form on windows: %s", absPath)
	}
	// Flip drive letter case
	altDrive := absPath
	if altDrive[0] >= 'a' && altDrive[0] <= 'z' {
		altDrive = string(absPath[0]-('a'-'A')) + absPath[1:]
	} else {
		altDrive = string(absPath[0]+('a'-'A')) + absPath[1:]
	}

	h1, err := fileHash(absPath)
	if err != nil {
		t.Fatalf("fileHash original: %v", err)
	}
	h2, err := fileHash(altDrive)
	if err != nil {
		t.Fatalf("fileHash alt: %v", err)
	}
	if h1 != h2 {
		t.Errorf("fileHash mismatch for case variants: %s vs %s", h1, h2)
	}

	p1, err := pathHash(absPath)
	if err != nil {
		t.Fatalf("pathHash original: %v", err)
	}
	p2, err := pathHash(altDrive)
	if err != nil {
		t.Fatalf("pathHash alt: %v", err)
	}
	if p1 != p2 {
		t.Errorf("pathHash mismatch for case variants: %s vs %s", p1, p2)
	}

	_ = wd // silence if unused later
}
