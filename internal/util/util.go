package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// ResolveHome resolves the home directory, preferring REAL_HOME for Docker usage
func ResolveHome(realHome string) string {
	if realHome != "" {
		return realHome
	}
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return h
}

// TitleCase converts the first character to uppercase
func TitleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// CurrentUser returns the current username
func CurrentUser() string {
	u, err := user.Current()
	if err != nil {
		return "user"
	}
	return u.Username
}

// Run executes a command with stdout/stderr connected
func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCapture executes a command and returns the combined output
func RunCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// EnsureFile creates a file if it doesn't exist
func EnsureFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	return f.Close()
}

// OK prints a green checkmark message
func OK(msg string) {
	fmt.Printf("\x1b[32m✓\x1b[0m %s\n", msg)
}

// Warn prints a yellow warning message
func Warn(msg string) {
	fmt.Printf("\x1b[33m•\x1b[0m %s\n", msg)
}

// Err prints a red error message
func Err(msg string) {
	fmt.Printf("\x1b[31m✗ %s\x1b[0m\n", msg)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsGlabAvailable checks if glab CLI is available
func IsGlabAvailable() bool {
	_, err := exec.LookPath("glab")
	return err == nil
}

// GetGitConfigHome returns the git config directory
func GetGitConfigHome(home string) string {
	return filepath.Join(home, ".gitconfig")
}
