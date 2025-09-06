package gitops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
)

// WritePerAccountGitConfig creates a per-account Git configuration
func WritePerAccountGitConfig(acc, name, email, keyPath, gitAlias, gitHost, folder, home string) error {
	cfgPath := filepath.Join(home, ".gitconfig-"+acc)

	config := fmt.Sprintf(`[user]
    name = %s
    email = %s
    signingkey = %s.pub
[gpg]
    format = ssh
[commit]
    gpgsign = true
[tag]
    gpgsign = true
[core]
    sshCommand = ssh -F ~/.ssh/config -i %s
[url "ssh://git@%s/"]
    insteadOf = https://%s/
    insteadOf = git@%s:
`, name, email, keyPath, keyPath, gitAlias, gitHost, gitHost)

	if err := os.WriteFile(cfgPath, []byte(config), 0o600); err != nil {
		return fmt.Errorf("failed to write Git config: %w", err)
	}
	util.OK("Wrote " + cfgPath)

	// Add includeIf binding
	if err := util.Run("git", "config", "--global", "--add",
		"includeIf.gitdir:"+trimTrailingSlash(folder)+"/.path", cfgPath); err != nil {
		return fmt.Errorf("failed to add includeIf binding: %w", err)
	}
	util.OK("Bound " + folder + " → " + filepath.Base(cfgPath))
	return nil
}

// trimTrailingSlash removes trailing slashes from a path
func trimTrailingSlash(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '/' || s[len(s)-1] == '\\') {
		s = s[:len(s)-1]
	}
	return s
}

// EnsureProjectFolder creates the project folder if it doesn't exist
func EnsureProjectFolder(folder string) error {
	if err := util.EnsureDir(folder); err != nil {
		return fmt.Errorf("failed to create project folder: %w", err)
	}
	util.OK("Project folder ready: " + folder)
	return nil
}

// CheckGitConfig verifies if Git is configured properly for an account
func CheckGitConfig(acc, folder, home string) error {
	cfgPath := filepath.Join(home, ".gitconfig-"+acc)

	// Check if account config exists
	if !util.FileExists(cfgPath) {
		return fmt.Errorf("Git config for account %s not found at %s", acc, cfgPath)
	}

	// Check if includeIf is configured
	output, err := util.RunCapture("git", "config", "--global", "--get-regexp", "includeIf")
	if err != nil {
		return fmt.Errorf("failed to check Git includeIf configuration: %w", err)
	}

	expectedPattern := "includeIf.gitdir:" + trimTrailingSlash(folder) + "/.path"
	if !strings.Contains(output, expectedPattern) {
		return fmt.Errorf("includeIf not configured for folder %s", folder)
	}

	util.OK("Git configuration validated for account: " + acc)
	return nil
}

// GetEffectiveGitConfig shows the effective Git configuration in a directory
func GetEffectiveGitConfig(folder string) error {
	// Change to the folder temporarily
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(folder); err != nil {
		return fmt.Errorf("failed to change to folder %s: %w", folder, err)
	}
	defer os.Chdir(originalDir)

	fmt.Printf("=== Effective Git config in %s ===\n", folder)
	fmt.Println("User name:")
	_ = util.Run("git", "config", "user.name")
	fmt.Println("User email:")
	_ = util.Run("git", "config", "user.email")
	fmt.Println("Signing key:")
	_ = util.Run("git", "config", "user.signingkey")
	fmt.Println("SSH command:")
	_ = util.Run("git", "config", "core.sshCommand")

	return nil
}
