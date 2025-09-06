package cmd

import (
	"fmt"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/sshops"
	"github.com/spf13/cobra"
)

var (
	testAlias string
	testRepo  string
)

func init() {
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test SSH auth to alias and optionally git ls-remote on a repo",
		Long: `Test SSH authentication and Git connectivity.

This command performs:
- SSH authentication test to the specified alias
- Optional Git ls-remote test on a repository
- Connection validation and troubleshooting

Examples:
  # Test SSH authentication only
  gvm-ssh test --alias gitlab-git

  # Test SSH and Git connectivity
  gvm-ssh test --alias gitlab-git --repo Global-Vision-Media/my-project

  # Test deploy key connectivity
  gvm-ssh test --alias ci-git --repo myorg/myproject`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if testAlias == "" {
				return fmt.Errorf("--alias is required")
			}

			// Test SSH authentication
			if err := sshops.TestSSHAuth(testAlias); err != nil {
				return fmt.Errorf("SSH auth test failed: %w", err)
			}

			// Optionally test Git remote access
			if testRepo != "" {
				fmt.Println()
				if err := sshops.TestGitRemote(testAlias, testRepo); err != nil {
					return fmt.Errorf("Git remote test failed: %w", err)
				}
			}

			return nil
		},
	}

	testCmd.Flags().StringVar(&testAlias, "alias", "", "SSH alias (e.g., gitlab-git)")
	testCmd.Flags().StringVar(&testRepo, "repo", "", "Namespace/repo for ls-remote test (e.g., Global-Vision-Media/project)")

	testCmd.MarkFlagRequired("alias")

	rootCmd.AddCommand(testCmd)
}
