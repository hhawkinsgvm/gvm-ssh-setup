package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/auth"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/gitops"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/sshops"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
)

var (
	acc        string
	gitAlias   string
	gitHost    string
	gitPort    int
	adminAlias string
	adminHost  string
	adminPort  int
	folder     string
	keyPath    string
	passphrase bool
	nameFlag   string
	emailFlag  string
	pushFlag   string // reserved for future (github|gitlab|none)
	
	// GitLab specific flags
	gitlabGroup  string
	uploadKey    bool
	skipAuth     bool
)

func init() {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Non-interactive setup of SSH aliases, key, and per-directory Git config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			home := util.ResolveHome(RealHome)
			if gitHost == "" || gitAlias == "" || acc == "" {
				return fmt.Errorf("--account, --git-alias, and --git-host are required")
			}
			
			// Check GitLab group membership unless skipped
			if !skipAuth && gitlabGroup != "" {
				util.Warn("Checking GitLab group membership...")
				if err := auth.CheckGitLabAccess(gitHost, gitlabGroup); err != nil {
					return fmt.Errorf("authorization failed: %w", err)
				}
				util.OK("GitLab access verified")
			}
			
			if folder == "" {
				folder = filepath.Join(home, "projects", acc) + string(os.PathSeparator)
			}
			if keyPath == "" {
				keyPath = filepath.Join(home, ".ssh", "id_ed25519_"+acc)
			}
			if nameFlag == "" {
				nameFlag = util.TitleCase(acc)
			}
			if emailFlag == "" {
				// Default to GVM email pattern, but allow override
				if strings.Contains(gitHost, "globalvision") {
					emailFlag = fmt.Sprintf("hud+%s@globalvision.com.au", acc)
				} else {
					emailFlag = fmt.Sprintf("%s@example.local", acc)
				}
			}

			// Ensure ~/.ssh/config includes config.d
			if err := sshops.EnsureInclude(home); err != nil {
				return err
			}

			// Generate/reuse key; add to agent
			if err := sshops.EnsureKey(keyPath, acc, passphrase); err != nil {
				return err
			}

			// Upload key to GitLab if requested
			if uploadKey && gitHost != "" {
				apiBase := "https://" + gitHost
				keyTitle := fmt.Sprintf("gvm-ssh-%s", acc)
				if err := sshops.AddSSHKeyToGitLab(apiBase, keyPath, keyTitle); err != nil {
					util.Warn(fmt.Sprintf("Failed to upload key to GitLab: %v", err))
					util.Warn("You can manually add the key in GitLab → Profile → SSH Keys")
				}
			}

			// Write git alias (user git)
			if err := sshops.WriteAlias(gitAlias, gitHost, gitPort, "git", keyPath, home); err != nil {
				return err
			}

			// Optional admin alias (user = current user)
			if adminAlias != "" && adminHost != "" {
				if err := sshops.WriteAlias(adminAlias, adminHost, adminPort, util.CurrentUser(), keyPath, home); err != nil {
					return err
				}
			}

			// Per-directory git identity + includeIf
			if err := gitops.WritePerAccountGitConfig(acc, nameFlag, emailFlag, keyPath, gitAlias, gitHost, folder, home); err != nil {
				return err
			}

			// Show fingerprints
			if err := sshops.ShowFingerprints(gitHost, gitPort, keyPath+".pub"); err != nil {
				return err
			}

			util.OK("Done")
			fmt.Printf("Folder: %s\n", folder)
			fmt.Printf("Remote example: git@%s:namespace/repo.git\n", gitAlias)
			
			if uploadKey {
				fmt.Printf("SSH key uploaded to GitLab: %s\n", gitHost)
			} else {
				fmt.Printf("To upload your SSH key manually, visit: https://%s/-/profile/keys\n", gitHost)
			}
			
			return nil
		},
	}

	setupCmd.Flags().StringVar(&acc, "account", "", "Account name (folder under ~/projects/)")
	setupCmd.Flags().StringVar(&gitAlias, "git-alias", "", "SSH alias for Git operations (e.g., gitlab-git)")
	setupCmd.Flags().StringVar(&gitHost, "git-host", "", "Hostname for Git operations")
	setupCmd.Flags().IntVar(&gitPort, "git-port", 22, "Port for Git operations")
	setupCmd.Flags().StringVar(&adminAlias, "admin-alias", "", "SSH alias for admin/shell access (optional)")
	setupCmd.Flags().StringVar(&adminHost, "admin-host", "", "Hostname/IP for admin access (optional)")
	setupCmd.Flags().IntVar(&adminPort, "admin-port", 22, "Port for admin access")
	setupCmd.Flags().StringVar(&folder, "folder", "", "Projects root folder (defaults to ~/projects/<account>/)")
	setupCmd.Flags().StringVar(&keyPath, "key", "", "Path to SSH key (defaults to ~/.ssh/id_ed25519_<account>)")
	setupCmd.Flags().BoolVar(&passphrase, "passphrase", false, "Protect the key with a passphrase")
	setupCmd.Flags().StringVar(&nameFlag, "name", "", "Git display name")
	setupCmd.Flags().StringVar(&emailFlag, "email", "", "Git email")
	setupCmd.Flags().StringVar(&pushFlag, "push", "none", "Reserved: github|gitlab|none")
	
	// GitLab specific flags
	setupCmd.Flags().StringVar(&gitlabGroup, "gitlab-group", "", "GitLab group to verify membership (required for authorization)")
	setupCmd.Flags().BoolVar(&uploadKey, "upload-key", false, "Upload SSH key to GitLab automatically")
	setupCmd.Flags().BoolVar(&skipAuth, "skip-auth", false, "Skip GitLab authorization check (for testing)")

	rootCmd.AddCommand(setupCmd)
}