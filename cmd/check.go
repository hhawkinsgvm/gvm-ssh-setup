package cmd

import (
	"fmt"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/sshops"
	"github.com/spf13/cobra"
)

var (
	checkAlias string
	checkHost  string
	checkPort  int
)

func init() {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Show effective SSH config for an alias and (optionally) host fingerprints",
		Long: `Check SSH configuration and connectivity.

This command displays:
- Effective SSH configuration for an alias
- Host key fingerprints for verification
- Connection test results

Examples:
  # Check SSH alias configuration
  gvm-ssh check --alias gitlab-git

  # Check configuration and scan host keys
  gvm-ssh check --alias gitlab-git --git-host gitlab.globalvision.com.au

  # Check with custom port
  gvm-ssh check --alias gitlab-git --git-host gitlab.globalvision.com.au --git-port 2122`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if checkAlias == "" {
				return fmt.Errorf("--alias is required")
			}

			// Show effective SSH configuration
			if err := sshops.ShowEffectiveConfig(checkAlias); err != nil {
				return fmt.Errorf("failed to show SSH config: %w", err)
			}

			// Optionally scan host keys
			if checkHost != "" {
				fmt.Println()
				if err := sshops.ScanHostKeys(checkHost, checkPort); err != nil {
					return fmt.Errorf("failed to scan host keys: %w", err)
				}
			}

			return nil
		},
	}

	checkCmd.Flags().StringVar(&checkAlias, "alias", "", "SSH alias to check")
	checkCmd.Flags().StringVar(&checkHost, "git-host", DefaultGitHost, "Hostname for fingerprint scan (optional)")
	checkCmd.Flags().IntVar(&checkPort, "git-port", DefaultGitPort, "Port for fingerprint scan (optional)")

	checkCmd.MarkFlagRequired("alias")

	rootCmd.AddCommand(checkCmd)
}
