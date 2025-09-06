package sshops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
)

// SSH directory and config paths
func sshDir(home string) string {
	return filepath.Join(home, ".ssh")
}

func confDir(home string) string {
	return filepath.Join(home, ".ssh", "config.d")
}

func mainConfig(home string) string {
	return filepath.Join(home, ".ssh", "config")
}

// EnsureInclude ensures that ~/.ssh/config includes the config.d directory
func EnsureInclude(home string) error {
	if err := util.EnsureDir(confDir(home)); err != nil {
		return err
	}
	if err := util.EnsureFile(mainConfig(home)); err != nil {
		return err
	}

	data, _ := os.ReadFile(mainConfig(home))
	includeDirective := "Include ~/.ssh/config.d/*.conf"

	if !strings.Contains(string(data), includeDirective) {
		f, err := os.OpenFile(mainConfig(home), os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = f.WriteString(includeDirective + "\n"); err != nil {
			return err
		}
		util.OK("Enabled config.d include in ~/.ssh/config")
	}
	return nil
}

// EnsureKey generates an SSH key if it doesn't exist and adds it to the agent
func EnsureKey(keyPath, comment string, passphrase bool) error {
	if util.FileExists(keyPath) {
		util.Warn("Key exists: " + keyPath)
	} else {
		args := []string{"-t", "ed25519", "-a", "100", "-C", comment, "-f", keyPath}
		if !passphrase {
			args = append(args, "-N", "")
		}
		if err := util.Run("ssh-keygen", args...); err != nil {
			return fmt.Errorf("failed to generate SSH key: %w", err)
		}
		util.OK("Generated key: " + keyPath)
	}

	// Try to add to ssh-agent
	_ = util.Run("bash", "-lc", `eval $(ssh-agent -s) >/dev/null 2>&1 || true`)
	_ = util.Run("ssh-add", keyPath)

	return nil
}

// WriteAlias creates an SSH alias configuration file
func WriteAlias(alias, host string, port int, user, keyPath, home string) error {
	if err := util.EnsureDir(confDir(home)); err != nil {
		return err
	}

	config := fmt.Sprintf(`Host %s
  HostName %s
  User %s
  Port %d
  IdentityFile %s
  IdentitiesOnly yes
`, alias, host, user, port, keyPath)

	configPath := filepath.Join(confDir(home), alias+".conf")
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		return fmt.Errorf("failed to write SSH alias config: %w", err)
	}

	if err := os.Chmod(configPath, 0o600); err != nil {
		return fmt.Errorf("failed to set permissions on SSH config: %w", err)
	}

	util.OK(fmt.Sprintf("Alias %s → %s@%s:%d", alias, user, host, port))
	return nil
}

// ShowFingerprints displays the fingerprint of the user's key and the host's keys
func ShowFingerprints(host string, port int, keyPubPath string) error {
	fmt.Println("=== Your key fingerprint ===")
	if err := util.Run("ssh-keygen", "-lf", keyPubPath); err != nil {
		util.Warn("Could not display key fingerprint")
	}

	fmt.Printf("=== %s:%d host keys (verify!) ===\n", host, port)
	cmd := fmt.Sprintf("ssh-keyscan -p %d %s 2>/dev/null | ssh-keygen -lf -", port, host)
	if err := util.Run("bash", "-lc", cmd); err != nil {
		util.Warn("Could not scan host keys")
	}

	return nil
}

// ShowEffectiveConfig shows the effective SSH configuration for an alias
func ShowEffectiveConfig(alias string) error {
	fmt.Printf("---- effective ssh -G %s ----\n", alias)
	return util.Run("ssh", "-G", alias)
}

// ScanHostKeys scans and displays host keys for verification
func ScanHostKeys(host string, port int) error {
	fmt.Printf("---- host fingerprints %s:%d ----\n", host, port)
	cmd := fmt.Sprintf("ssh-keyscan -p %d %s 2>/dev/null | ssh-keygen -lf -", port, host)
	return util.Run("bash", "-lc", cmd)
}

// TestSSHAuth tests SSH authentication to a host alias
func TestSSHAuth(alias string) error {
	util.Warn("SSH auth test (BatchMode)…")
	out, _ := util.RunCapture("ssh", "-T", "-o", "BatchMode=yes", "git@"+alias)
	fmt.Print(out)
	util.OK("SSH test executed (some servers close after banner; this can still be OK)")
	return nil
}

// TestGitRemote tests git ls-remote to verify Git connectivity
func TestGitRemote(alias, nsrepo string) error {
	util.Warn("git ls-remote test…")
	cmd := fmt.Sprintf(`GIT_SSH_COMMAND="ssh -F ~/.ssh/config" git ls-remote "git@%s:%s.git"`, alias, nsrepo)
	if err := util.Run("bash", "-lc", cmd); err != nil {
		return fmt.Errorf("git remote test failed for %s: %w", nsrepo, err)
	}
	util.OK("Git remote reachable: " + nsrepo)
	return nil
}

// ReadPublicKey reads the content of a public key file
func ReadPublicKey(keyPath string) (string, error) {
	pubKeyPath := keyPath + ".pub"
	content, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key %s: %w", pubKeyPath, err)
	}
	return strings.TrimSpace(string(content)), nil
}

// GetKeyFingerprint returns the fingerprint of an SSH key
func GetKeyFingerprint(keyPath string) (string, error) {
	pubKeyPath := keyPath + ".pub"
	output, err := util.RunCapture("ssh-keygen", "-lf", pubKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to get key fingerprint: %w", err)
	}
	return strings.TrimSpace(output), nil
}
