package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	wizardCmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive guided setup for GVM GitLab CE environment",
		Long: `Interactive wizard for setting up SSH and Git configuration.

This wizard walks you through the complete setup process with
sensible defaults for Global Vision Media's GitLab CE environment.

The wizard will:
- Guide you through account setup
- Configure SSH keys and aliases
- Set up per-directory Git identity
- Optionally upload keys to GitLab
- Test connectivity

Pre-configured defaults:
- Git Host: gitlab.globalvision.com.au
- Git Port: 2122
- Admin Host: 203.32.94.10
- Admin Port: 2122
- Email Domain: @globalvision.com.au`,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := bufio.NewReader(os.Stdin)
			home := util.ResolveHome(RealHome)

			fmt.Println("🚀 GVM SSH Setup Wizard")
			fmt.Println("======================")
			fmt.Println()

			// Check for GitLab token
			hasGitLabToken := false
			if token := os.Getenv("GITLAB_TOKEN"); token != "" {
				hasGitLabToken = true
				util.OK("GitLab token found in environment")
			} else if util.IsGlabAvailable() {
				util.Warn("glab CLI available but no GITLAB_TOKEN set")
				fmt.Println("Run 'glab auth login' or set GITLAB_TOKEN for key upload")
			} else {
				util.Warn("No GitLab authentication available")
				fmt.Println("Install glab CLI or set GITLAB_TOKEN for key upload")
			}
			fmt.Println()

			// Basic configuration
			acc := ask(r, "Account name", "gvm")
			gitAlias := ask(r, "Git SSH alias", acc+"-git")
			gitHost := ask(r, "Git hostname", DefaultGitHost)
			gitPortStr := ask(r, "Git port", fmt.Sprintf("%d", DefaultGitPort))
			gitPort := atoi(gitPortStr)

			// Admin configuration
			fmt.Println()
			useAdmin := askBool(r, "Configure admin/shell access?", true)
			var adminAlias, adminHost string
			adminPort := DefaultAdminPort
			if useAdmin {
				adminAlias = ask(r, "Admin SSH alias", acc+"-host")
				adminHost = ask(r, "Admin hostname", DefaultAdminHost)
				adminPortStr := ask(r, "Admin port", fmt.Sprintf("%d", DefaultAdminPort))
				adminPort = atoi(adminPortStr)
			}

			// Project and identity configuration
			fmt.Println()
			folder := ask(r, "Projects folder", filepath.Join(home, "projects", acc)+"/")
			name := ask(r, "Git display name", fmt.Sprintf("GVM %s", util.TitleCase(acc)))
			email := ask(r, "Git email", fmt.Sprintf("%s@globalvision.com.au", acc))
			keyPath := ask(r, "SSH key path", filepath.Join(home, ".ssh", "id_ed25519_"+acc))

			// Security options
			fmt.Println()
			usePassphrase := askBool(r, "Protect SSH key with passphrase?", true)

			// GitLab integration
			uploadKey := false
			if hasGitLabToken {
				fmt.Println()
				uploadKey = askBool(r, "Upload SSH key to GitLab?", true)
				if uploadKey {
					orgGroup := ask(r, "GitLab group for access verification", DefaultOrgGroup)
					// Store the org group for the setup command
					DefaultOrgGroup = orgGroup
				}
			}

			// Confirmation
			fmt.Println()
			fmt.Println("📋 Configuration Summary:")
			fmt.Printf("  Account: %s\n", acc)
			fmt.Printf("  Git Alias: %s (%s:%d)\n", gitAlias, gitHost, gitPort)
			if useAdmin {
				fmt.Printf("  Admin Alias: %s (%s:%d)\n", adminAlias, adminHost, adminPort)
			}
			fmt.Printf("  Project Folder: %s\n", folder)
			fmt.Printf("  Git Identity: %s <%s>\n", name, email)
			fmt.Printf("  SSH Key: %s\n", keyPath)
			fmt.Printf("  Passphrase Protected: %v\n", usePassphrase)
			if uploadKey {
				fmt.Printf("  Upload to GitLab: %v\n", uploadKey)
			}
			fmt.Println()

			if !askBool(r, "Proceed with setup?", true) {
				fmt.Println("Setup cancelled.")
				return nil
			}

			// Build arguments for setup command
			setupArgs := []string{
				"setup",
				"--account", acc,
				"--git-alias", gitAlias,
				"--git-host", gitHost,
				"--git-port", fmt.Sprintf("%d", gitPort),
				"--folder", folder,
				"--name", name,
				"--email", email,
				"--key", keyPath,
			}

			if usePassphrase {
				setupArgs = append(setupArgs, "--passphrase")
			}

			if useAdmin && adminAlias != "" && adminHost != "" {
				setupArgs = append(setupArgs,
					"--admin-alias", adminAlias,
					"--admin-host", adminHost,
					"--admin-port", fmt.Sprintf("%d", adminPort))
			}

			if uploadKey {
				setupArgs = append(setupArgs, "--upload-key")
			}

			// Execute setup command
			fmt.Println()
			fmt.Println("🔧 Executing setup...")

			// Temporarily modify os.Args to call setup command
			originalArgs := os.Args
			os.Args = append([]string{os.Args[0]}, setupArgs...)
			defer func() { os.Args = originalArgs }()

			return rootCmd.Execute()
		},
	}

	rootCmd.AddCommand(wizardCmd)
}

// ask prompts the user with a question and default value
func ask(r *bufio.Reader, prompt, defaultVal string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultVal)
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

// askBool prompts the user for a yes/no answer
func askBool(r *bufio.Reader, prompt string, defaultVal bool) bool {
	defaultStr := "n"
	if defaultVal {
		defaultStr = "y"
	}

	response := ask(r, prompt+" (y/n)", defaultStr)
	response = strings.ToLower(response)
	return response == "y" || response == "yes"
}

// atoi converts string to int, returning 0 on error
func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
