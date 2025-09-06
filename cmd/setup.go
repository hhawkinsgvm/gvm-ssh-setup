package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/auth"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/gitops"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/sshops"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
	"github.com/spf13/cobra"
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
	uploadKey  bool
	skipAuth   bool
	deployKey  bool
	targetRepo string
	orgGroup   string
)

func init() {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Non-interactive setup of SSH aliases, key, and per-directory Git config",
		Long: `Setup SSH and Git configuration for a specific account.

This command creates:
- SSH key pair (if not exists)
- SSH alias configuration
- Per-directory Git identity
- Optional GitLab key upload
- Optional deploy key for CI/servers

Examples:
  # Basic developer setup
  gvm-ssh setup --account myaccount --git-alias myaccount-git

  # Full setup with admin access
  gvm-ssh setup --account gvm --git-alias gitlab-git \
    --admin-alias gitlab-host --upload-key

  # Deploy key for CI/server
  gvm-ssh setup --account ci --git-alias ci-git --deploy-key \
    --target-repo myorg/myproject --skip-auth`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			home := util.ResolveHome(RealHome)

			// Validate required flags
			if gitHost == "" || gitAlias == "" || acc == "" {
				return fmt.Errorf("--account, --git-alias, and --git-host are required")
			}

			// Set defaults
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
				emailFlag = fmt.Sprintf("%s@globalvision.com.au", acc)
			}
			if orgGroup == "" {
				orgGroup = DefaultOrgGroup
			}

			// GitLab authentication and authorization check
			var gitlabConfig *auth.GitLabConfig
			if !skipAuth {
				token := auth.GetGitLabToken()
				if token == "" {
					return fmt.Errorf("GitLab token required. Set GITLAB_TOKEN or run 'glab auth login'")
				}

				var err error
				gitlabConfig, err = auth.NewGitLabConfig(GitLabBaseURL, token)
				if err != nil {
					return fmt.Errorf("failed to configure GitLab client: %w", err)
				}

				// Check group membership (unless it's a deploy key setup)
				if !deployKey {
					isMember, err := gitlabConfig.CheckGroupMembership(orgGroup)
					if err != nil {
						return fmt.Errorf("failed to check group membership: %w", err)
					}
					if !isMember {
						return fmt.Errorf("access denied: you are not a member of %s", orgGroup)
					}
					util.OK(fmt.Sprintf("Verified membership in %s", orgGroup))
				}
			}

			// Ensure ~/.ssh/config includes config.d
			if err := sshops.EnsureInclude(home); err != nil {
				return fmt.Errorf("failed to setup SSH config include: %w", err)
			}

			// Generate/reuse key; add to agent
			if err := sshops.EnsureKey(keyPath, acc, passphrase); err != nil {
				return fmt.Errorf("failed to ensure SSH key: %w", err)
			}

			// Write git alias (user git)
			if err := sshops.WriteAlias(gitAlias, gitHost, gitPort, "git", keyPath, home); err != nil {
				return fmt.Errorf("failed to write Git SSH alias: %w", err)
			}

			// Optional admin alias (user = current user)
			if adminAlias != "" && adminHost != "" {
				if err := sshops.WriteAlias(adminAlias, adminHost, adminPort, util.CurrentUser(), keyPath, home); err != nil {
					return fmt.Errorf("failed to write admin SSH alias: %w", err)
				}
			}

			// Per-directory git identity + includeIf (skip for deploy keys)
			if !deployKey {
				if err := gitops.EnsureProjectFolder(folder); err != nil {
					return fmt.Errorf("failed to ensure project folder: %w", err)
				}

				if err := gitops.WritePerAccountGitConfig(acc, nameFlag, emailFlag, keyPath, gitAlias, gitHost, folder, home); err != nil {
					return fmt.Errorf("failed to write Git config: %w", err)
				}
			}

			// Upload SSH key to GitLab
			if uploadKey && gitlabConfig != nil {
				pubKeyContent, err := sshops.ReadPublicKey(keyPath)
				if err != nil {
					return fmt.Errorf("failed to read public key: %w", err)
				}

				keyTitle := fmt.Sprintf("gvm-ssh-%s", acc)
				if deployKey {
					keyTitle = fmt.Sprintf("deploy-%s", acc)
				}

				if err := gitlabConfig.AddSSHKey(keyTitle, pubKeyContent); err != nil {
					util.Warn(fmt.Sprintf("Failed to upload SSH key: %v", err))
					util.Warn("You can manually add the key in GitLab > Preferences > SSH Keys")
				} else {
					util.OK("SSH key uploaded to GitLab")
				}
			}

			// Show fingerprints for verification
			if err := sshops.ShowFingerprints(gitHost, gitPort, keyPath+".pub"); err != nil {
				util.Warn("Could not display fingerprints")
			}

			util.OK("Setup completed successfully!")

			// Display usage information
			fmt.Printf("\nSetup Summary:\n")
			fmt.Printf("Account: %s\n", acc)
			fmt.Printf("SSH Alias: %s\n", gitAlias)
			if !deployKey {
				fmt.Printf("Project Folder: %s\n", folder)
				fmt.Printf("Git Identity: %s <%s>\n", nameFlag, emailFlag)
			}
			fmt.Printf("Key Path: %s\n", keyPath)
			if adminAlias != "" {
				fmt.Printf("Admin Alias: %s\n", adminAlias)
			}

			fmt.Printf("\nUsage Examples:\n")
			fmt.Printf("  git clone git@%s:namespace/repo.git\n", gitAlias)
			if adminAlias != "" {
				fmt.Printf("  ssh %s\n", adminAlias)
			}

			return nil
		},
	}

	// Flags
	setupCmd.Flags().StringVar(&acc, "account", "", "Account name (folder under ~/projects/)")
	setupCmd.Flags().StringVar(&gitAlias, "git-alias", "", "SSH alias for Git operations (e.g., gitlab-git)")
	setupCmd.Flags().StringVar(&gitHost, "git-host", DefaultGitHost, "Hostname for Git operations")
	setupCmd.Flags().IntVar(&gitPort, "git-port", DefaultGitPort, "Port for Git operations")
	setupCmd.Flags().StringVar(&adminAlias, "admin-alias", "", "SSH alias for admin/shell access (optional)")
	setupCmd.Flags().StringVar(&adminHost, "admin-host", DefaultAdminHost, "Hostname/IP for admin access (optional)")
	setupCmd.Flags().IntVar(&adminPort, "admin-port", DefaultAdminPort, "Port for admin access")
	setupCmd.Flags().StringVar(&folder, "folder", "", "Projects root folder (defaults to ~/projects/<account>/)")
	setupCmd.Flags().StringVar(&keyPath, "key", "", "Path to SSH key (defaults to ~/.ssh/id_ed25519_<account>)")
	setupCmd.Flags().BoolVar(&passphrase, "passphrase", false, "Protect the key with a passphrase")
	setupCmd.Flags().StringVar(&nameFlag, "name", "", "Git display name (defaults to title-cased account)")
	setupCmd.Flags().StringVar(&emailFlag, "email", "", "Git email (defaults to <account>@globalvision.com.au)")
	setupCmd.Flags().StringVar(&pushFlag, "push", "none", "Reserved: github|gitlab|none")
	setupCmd.Flags().BoolVar(&uploadKey, "upload-key", false, "Upload SSH key to GitLab account")
	setupCmd.Flags().BoolVar(&skipAuth, "skip-auth", false, "Skip GitLab authentication and group membership check")
	setupCmd.Flags().BoolVar(&deployKey, "deploy-key", false, "Setup as deploy key (for CI/servers)")
	setupCmd.Flags().StringVar(&targetRepo, "target-repo", "", "Target repository for deploy key (org/repo)")
	setupCmd.Flags().StringVar(&orgGroup, "org-group", DefaultOrgGroup, "GitLab group to check membership against")

	// Mark required flags
	setupCmd.MarkFlagRequired("account")
	setupCmd.MarkFlagRequired("git-alias")

	rootCmd.AddCommand(setupCmd)
}
