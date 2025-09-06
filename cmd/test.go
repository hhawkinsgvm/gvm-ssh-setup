package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/sshops"
)

var (
	testAlias string
	testRepo  string
)

func init() {
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test SSH auth to alias and optionally git ls-remote on a repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if testAlias == "" {
				return fmt.Errorf("--alias is required")
			}
			if err := sshops.TestSSHAuth(testAlias); err != nil {
				return err
			}
			if testRepo != "" {
				return sshops.TestGitRemote(testAlias, testRepo)
			}
			return nil
		},
	}
	testCmd.Flags().StringVar(&testAlias, "alias", "", "SSH alias (e.g., gitlab-git)")
	testCmd.Flags().StringVar(&testRepo, "repo", "", "Namespace/repo for ls-remote (e.g., Org/Repo)")
	rootCmd.AddCommand(testCmd)
}