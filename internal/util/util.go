package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

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

func TitleCase(s string) string {
	if s == "" { return s }
	return strings.ToUpper(s[:1]) + s[1:]
}

func CurrentUser() string {
	u, err := user.Current()
	if err != nil {
		return "user"
	}
	return u.Username
}

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func EnsureDir(path string) error { return os.MkdirAll(path, 0o755) }

func EnsureFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil { return err }
	return f.Close()
}

func OK(msg string)   { fmt.Printf("\x1b[32m✓\x1b[0m %s\n", msg) }
func Warn(msg string) { fmt.Printf("\x1b[33m•\x1b[0m %s\n", msg) }
func Err(msg string)  { fmt.Printf("\x1b[31m✗ %s\x1b[0m\n", msg) }