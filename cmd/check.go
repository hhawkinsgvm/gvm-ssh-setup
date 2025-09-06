package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/sshops"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if checkAlias == "" {
				return fmt.Errorf("--alias is required")
			}
			if err := sshops.ShowEffectiveConfig(checkAlias); err != nil {
				return err
			}
			if checkHost != "" {
				if err := sshops.ScanHostKeys(checkHost, checkPort); err != nil {
					return err
				}
			}
			return nil
		},
	}
	checkCmd.Flags().StringVar(&checkAlias, "alias", "", "SSH alias to check")
	checkCmd.Flags().StringVar(&checkHost, "git-host", "", "Hostname for fingerprint scan (optional)")
	checkCmd.Flags().IntVar(&checkPort, "git-port", 22, "Port for fingerprint scan (optional)")
	rootCmd.AddCommand(checkCmd)
}