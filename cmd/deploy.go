package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/auth"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/sshops"
	"github.com/hhawkinsgvm/gvm-ssh-setup/internal/util"
)

var (
	deployKeyPath     string
	deployKeyTitle    string
	deployProject     string
	deployReadOnly    bool
	deployGitlabGroup string
)

func init() {
	deployCmd := &cobra.Command{
		Use:   "deploy-key",
		Short: "Generate and manage deploy keys for servers/CI (admin only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			home := util.ResolveHome(RealHome)
			
			if deployProject == "" {
				return fmt.Errorf("--project is required")
			}
			
			// Check GitLab group membership for admin operations
			if deployGitlabGroup != "" {
				util.Warn("Checking GitLab admin access...")
				if err := auth.CheckGitLabAccess("gitlab.globalvision.com.au", deployGitlabGroup); err != nil {
					return fmt.Errorf("admin authorization failed: %w", err)
				}
				util.OK("GitLab admin access verified")
			}
			
			if deployKeyPath == "" {
				deployKeyPath = filepath.Join(home, ".ssh", "deploy_"+deployProject)
			}
			if deployKeyTitle == "" {
				deployKeyTitle = fmt.Sprintf("Deploy key for %s", deployProject)
			}
			
			// Generate deploy key
			if err := sshops.EnsureKey(deployKeyPath, deployKeyTitle, false); err != nil {
				return fmt.Errorf("failed to generate deploy key: %w", err)
			}
			
			util.OK("Deploy key generated: " + deployKeyPath)
			
			// Read and display the public key
			pubKey, err := os.ReadFile(deployKeyPath + ".pub")
			if err != nil {
				return fmt.Errorf("failed to read public key: %w", err)
			}
			
			fmt.Println("=== Deploy Key (add to GitLab project) ===")
			fmt.Printf("Title: %s\n", deployKeyTitle)
			fmt.Printf("Read-only: %t\n", deployReadOnly)
			fmt.Println("Key:")
			fmt.Print(string(pubKey))
			fmt.Println()
			
			fmt.Printf("Add this key manually in GitLab:\n")
			fmt.Printf("Project → Settings → Repository → Deploy Keys\n")
			fmt.Printf("Private key location: %s\n", deployKeyPath)
			
			return nil
		},
	}
	
	deployCmd.Flags().StringVar(&deployKeyPath, "key", "", "Path for deploy key (defaults to ~/.ssh/deploy_<project>)")
	deployCmd.Flags().StringVar(&deployKeyTitle, "title", "", "Title for the deploy key")
	deployCmd.Flags().StringVar(&deployProject, "project", "", "Project name for the deploy key (required)")
	deployCmd.Flags().BoolVar(&deployReadOnly, "read-only", true, "Create read-only deploy key")
	deployCmd.Flags().StringVar(&deployGitlabGroup, "gitlab-group", "", "GitLab group to verify admin membership")
	
	rootCmd.AddCommand(deployCmd)
}