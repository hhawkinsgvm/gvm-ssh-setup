package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// RealHome lets the app write to the host home when running in Docker
	RealHome = os.Getenv("REAL_HOME")

	// GitLab configuration
	GitLabBaseURL = os.Getenv("GITLAB_BASE_URL")
	GitLabToken   = os.Getenv("GITLAB_TOKEN")

	// Default values for GVM organization
	DefaultGitHost   = "gitlab.globalvision.com.au"
	DefaultGitPort   = 2122
	DefaultAdminHost = "203.32.94.10"
	DefaultAdminPort = 2122
	DefaultOrgGroup  = "global-vision-media" // Change this to your actual GitLab group path
)

var rootCmd = &cobra.Command{
	Use:   "gvm-ssh",
	Short: "GVM SSH/Git per-account setup tool for GitLab CE",
	Long: `GVM SSH Setup Tool

A comprehensive tool for setting up SSH aliases, keys, and per-directory Git identity
specifically optimized for GitLab CE environments.

Features:
- SSH key generation and management
- GitLab CE integration for key upload
- Per-account Git configuration with includeIf
- Deploy key support for servers/CI
- Interactive wizard and non-interactive setup modes
- Connectivity testing and configuration validation

This tool is designed for Global Vision Media's GitLab CE environment with
sensible defaults for common workflows.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip auth check for help and non-setup commands
		if cmd.Name() == "help" || cmd.Name() == "check" || cmd.Parent() == nil {
			return nil
		}

		// For setup and wizard commands, we'll validate GitLab access
		return nil
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Set default GitLab base URL if not specified
	if GitLabBaseURL == "" {
		GitLabBaseURL = "https://" + DefaultGitHost
	}
}
