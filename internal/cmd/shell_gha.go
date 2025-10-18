package cmd

import (
	"crypto/rand"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type gha struct{}

// GitHubActions shell instance
var GitHubActions Shell = gha{}

var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func (sh gha) Hook() (string, error) {
	return "", fmt.Errorf("Hook not implemented for GitHub Actions shell")
}

func (sh gha) Export(e ShellExport) (string, error) {
	var b strings.Builder
	for key, value := range e {
		if !validKeyPattern.MatchString(key) {
			// Skip invalid environment variable keys
			fmt.Fprintf(os.Stderr, "direnv: Skipping invalid environment variable key: %s\n", key)
			continue
		}
		if value == nil {
			sh.unset(&b, key)
		} else {
			if err := sh.export(&b, key, *value); err != nil {
				return "", err
			}
		}
	}
	return b.String(), nil
}

func (sh gha) Dump(env Env) (string, error) {
	var b strings.Builder

	for key, value := range env {
		if !validKeyPattern.MatchString(key) {
			// Skip invalid environment variable keys
			fmt.Fprintf(os.Stderr, "direnv: Skipping invalid environment variable key: %s\n", key)
			continue
		}
		if err := sh.export(&b, key, value); err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func (sh gha) export(b *strings.Builder, key, value string) error {
	// Generate a random delimiter
	delimiter := sh.generateDelimiter()

	// Check if key or value contains delimiter (should be extremely rare)
	if strings.Contains(key, delimiter) || strings.Contains(value, delimiter) {
		// Log the collision and regenerate delimiter
		fmt.Fprintf(os.Stderr, "direnv: Delimiter collision detected for key %s, regenerating delimiter\n", key)
		delimiter = sh.generateDelimiter()

		// If still colliding (astronomically unlikely), error out
		if strings.Contains(key, delimiter) || strings.Contains(value, delimiter) {
			return fmt.Errorf("delimiter collision after regeneration for key %s", key)
		}
	}

	b.WriteString(key)
	b.WriteString("<<")
	b.WriteString(delimiter)
	b.WriteByte('\n')
	b.WriteString(value)
	b.WriteByte('\n')
	b.WriteString(delimiter)
	b.WriteByte('\n')
	return nil
}

func (sh gha) unset(_ *strings.Builder, _ string) {
	// Don't do anything. > $GITHUB_ENV will overwrite the existing env.
}

func (sh gha) generateDelimiter() string {
	// Generate random bytes for delimiter
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to timestamp-based delimiter
		return fmt.Sprintf("ghadelimiter_%d", time.Now().UnixNano())
	}

	// Convert to hex string
	return fmt.Sprintf("ghadelimiter_%x", randomBytes)
}
