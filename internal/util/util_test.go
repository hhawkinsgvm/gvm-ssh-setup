package util

import (
	"os"
	"testing"
)

func TestResolveHome(t *testing.T) {
	// Test with REAL_HOME set
	realHome := "/test/home"
	result := ResolveHome(realHome)
	if result != realHome {
		t.Errorf("Expected %s, got %s", realHome, result)
	}
	
	// Test with empty REAL_HOME (should use system home)
	result = ResolveHome("")
	if result == "" {
		t.Error("Expected non-empty home directory")
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"HELLO", "HELLO"},
		{"", ""},
		{"h", "H"},
	}
	
	for _, test := range tests {
		result := TitleCase(test.input)
		if result != test.expected {
			t.Errorf("TitleCase(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestCurrentUser(t *testing.T) {
	user := CurrentUser()
	if user == "" {
		t.Error("Expected non-empty username")
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := "/tmp/test-ensure-dir"
	defer os.RemoveAll(tmpDir)
	
	err := EnsureDir(tmpDir)
	if err != nil {
		t.Errorf("EnsureDir failed: %v", err)
	}
	
	// Check directory exists
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

func TestEnsureFile(t *testing.T) {
	tmpFile := "/tmp/test-ensure-file"
	defer os.Remove(tmpFile)
	
	err := EnsureFile(tmpFile)
	if err != nil {
		t.Errorf("EnsureFile failed: %v", err)
	}
	
	// Check file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("File was not created")
	}
}

func TestRunCapture(t *testing.T) {
	// Test simple command that should work on any system
	output, err := RunCapture("echo", "hello")
	if err != nil {
		t.Errorf("RunCapture failed: %v", err)
	}
	
	if output != "hello\n" {
		t.Errorf("Expected 'hello\\n', got %q", output)
	}
}