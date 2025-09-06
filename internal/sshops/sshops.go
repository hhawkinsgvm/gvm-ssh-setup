package sshops

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
)

func sshDir(home string) string { return filepath.Join(home, ".ssh") }
func confDir(home string) string { return filepath.Join(home, ".ssh", "config.d") }
func mainConfig(home string) string { return filepath.Join(home, ".ssh", "config") }

func EnsureInclude(home string) error {
	if err := util.EnsureDir(confDir(home)); err != nil { return err }
	if err := util.EnsureFile(mainConfig(home)); err != nil { return err }
	data, _ := os.ReadFile(mainConfig(home))
	if !strings.Contains(string(data), "Include ~/.ssh/config.d/*.conf") {
		f, err := os.OpenFile(mainConfig(home), os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil { return err }
		defer f.Close()
		if _, err = f.WriteString("Include ~/.ssh/config.d/*.conf\n"); err != nil { return err }
		util.OK("Enabled config.d include in ~/.ssh/config")
	}
	return nil
}

func EnsureKey(keyPath, comment string, passphrase bool) error {
	if _, err := os.Stat(keyPath); err == nil {
		util.Warn("Key exists: " + keyPath)
	} else {
		args := []string{"-t", "ed25519", "-a", "100", "-C", comment, "-f", keyPath}
		if !passphrase {
			args = append(args, "-N", "")
		}
		if err := util.Run("ssh-keygen", args...); err != nil {
			return err
		}
		util.OK("Generated key: " + keyPath)
	}
	_ = util.Run("bash", "-lc", `eval $(ssh-agent -s) >/dev/null 2>&1 || true`)
	_ = util.Run("ssh-add", keyPath)
	return nil
}

func WriteAlias(alias, host string, port int, user, keyPath, home string) error {
	if err := util.EnsureDir(confDir(home)); err != nil { return err }
	text := fmt.Sprintf(`Host %s
  HostName %s
  User %s
  Port %d
  IdentityFile %s
  IdentitiesOnly yes
`, alias, host, user, port, keyPath)
	path := filepath.Join(confDir(home), alias+".conf")
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil { return err }
	if err := os.Chmod(path, 0o600); err != nil { return err }
	util.OK(fmt.Sprintf("Alias %s → %s@%s:%d", alias, user, host, port))
	return nil
}

func ShowFingerprints(host string, port int, keyPub string) error {
	fmt.Println("=== Your key fingerprint ===")
	_ = util.Run("ssh-keygen", "-lf", keyPub)
	fmt.Printf("=== %s:%d host keys (verify!) ===\n", host, port)
	_ = util.Run("bash", "-lc", fmt.Sprintf("ssh-keyscan -p %d %s 2>/dev/null | ssh-keygen -lf -", port, host))
	return nil
}

func ShowEffectiveConfig(alias string) error {
	fmt.Printf("---- effective ssh -G %s ----\n", alias)
	return util.Run("ssh", "-G", alias)
}

func ScanHostKeys(host string, port int) error {
	fmt.Printf("---- host fingerprints %s:%d ----\n", host, port)
	return util.Run("bash", "-lc", fmt.Sprintf("ssh-keyscan -p %d %s 2>/dev/null | ssh-keygen -lf -", port, host))
}

func TestSSHAuth(alias string) error {
	util.Warn("SSH auth test (BatchMode)…")
	out, _ := util.RunCapture("ssh", "-T", "-o", "BatchMode=yes", "git@"+alias)
	fmt.Print(out)
	util.OK("SSH test executed (some servers close after banner; this can still be OK)")
	return nil
}

func TestGitRemote(alias, nsrepo string) error {
	util.Warn("git ls-remote test…")
	cmd := fmt.Sprintf(`GIT_SSH_COMMAND="ssh -F ~/.ssh/config" git ls-remote "git@%s:%s.git"`, alias, nsrepo)
	if err := util.Run("bash", "-lc", cmd); err != nil {
		return fmt.Errorf("git remote test failed for %s: %w", nsrepo, err)
	}
	util.OK("Git remote reachable: " + nsrepo)
	return nil
}

// GitLabSSHKey represents an SSH key in GitLab API
type GitLabSSHKey struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Key   string `json:"key"`
}

// AddSSHKeyToGitLab uploads an SSH public key to GitLab
func AddSSHKeyToGitLab(apiBase, keyPath, title string) error {
	// First try using glab CLI if available
	if err := util.Run("glab", "version"); err == nil {
		util.Warn("Using glab CLI to add SSH key...")
		if err := util.Run("glab", "ssh-key", "add", keyPath+".pub", "-t", title); err != nil {
			util.Warn("glab failed, trying API directly...")
		} else {
			util.OK("SSH key added via glab CLI")
			return nil
		}
	}

	// Fallback to direct API call
	tok := os.Getenv("GITLAB_TOKEN")
	if tok == "" {
		return fmt.Errorf("GITLAB_TOKEN not set and glab not available; cannot upload key")
	}

	// Read the public key
	keyContent, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	// Prepare the request payload
	payload := map[string]string{
		"title": title,
		"key":   string(keyContent),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Make the API request
	req, err := http.NewRequest("POST", apiBase+"/api/v4/user/keys", strings.NewReader(string(payloadBytes)))
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", tok)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to upload key: HTTP %d", resp.StatusCode)
	}

	util.OK("SSH key uploaded to GitLab")
	return nil
}