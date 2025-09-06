package gitops

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
)

func WritePerAccountGitConfig(acc, name, email, keyPath, gitAlias, gitHost, folder, home string) error {
	cfgPath := filepath.Join(home, ".gitconfig-"+acc)
	cfg := fmt.Sprintf(`[user]
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

	if err := os.WriteFile(cfgPath, []byte(cfg), 0o600); err != nil {
		return err
	}
	util.OK("Wrote " + cfgPath)

	// includeIf binding
	if err := util.Run("git", "config", "--global", "--add",
		"includeIf.gitdir:"+trimTrailingSlash(folder)+"/.path", cfgPath); err != nil {
		return err
	}
	util.OK("Bound " + folder + " → " + filepath.Base(cfgPath))
	return nil
}

func trimTrailingSlash(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '/' || s[len(s)-1] == '\\') {
		s = s[:len(s)-1]
	}
	return s
}