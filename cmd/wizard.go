package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	wizardCmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive guided setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := bufio.NewReader(os.Stdin)
			home := os.Getenv("REAL_HOME")
			if home == "" {
				h, err := os.UserHomeDir()
				if err != nil { return err }
				home = h
			}

			fmt.Println("=== GVM SSH Setup Wizard ===")
			fmt.Println("This will configure SSH keys and Git identity for GitLab access.")
			fmt.Println()

			acc := ask(r, "Account name", "gvm")
			gitAlias := ask(r, "Git alias (git ops)", acc+"-git")
			gitHost := ask(r, "Git host", "gitlab.globalvision.com.au")
			gitPort := atoi(ask(r, "Git port", "2122"))

			// GitLab specific
			gitlabGroup := ask(r, "GitLab group (for access control)", "global-vision-media")
			uploadKey := strings.ToLower(ask(r, "Upload SSH key to GitLab automatically? (y/n)", "y"))

			useAdmin := strings.ToLower(ask(r, "Also add admin/shell alias? (y/n)", "y"))
			var adminAlias, adminHost string
			adminPort := 22
			if useAdmin == "y" || useAdmin == "yes" {
				adminAlias = ask(r, "Admin alias", acc+"-host")
				adminHost  = ask(r, "Admin host", "203.32.94.10")
				adminPort  = atoi(ask(r, "Admin port", "2122"))
			}

			folder := ask(r, "Projects folder", filepath.Join(home, "projects", acc)+"/")
			name   := ask(r, "Git display name", "Hud Hawkins ("+acc+")")
			email  := ask(r, "Git email", "hud+"+acc+"@globalvision.com.au")
			key    := ask(r, "Key path", filepath.Join(home, ".ssh", "id_ed25519_"+acc))
			pp     := strings.ToLower(ask(r, "Protect key with passphrase? (y/n)", "y"))
			passphrase := (pp == "y" || pp == "yes")

			fmt.Println()
			fmt.Println("=== Configuration Summary ===")
			fmt.Printf("Account: %s\n", acc)
			fmt.Printf("Git Alias: %s\n", gitAlias)
			fmt.Printf("Git Host: %s:%d\n", gitHost, gitPort)
			fmt.Printf("GitLab Group: %s\n", gitlabGroup)
			fmt.Printf("Upload Key: %s\n", uploadKey)
			fmt.Printf("Projects Folder: %s\n", folder)
			fmt.Printf("Git Name: %s\n", name)
			fmt.Printf("Git Email: %s\n", email)
			fmt.Printf("SSH Key: %s\n", key)
			fmt.Printf("Passphrase: %t\n", passphrase)
			fmt.Println()

			confirm := strings.ToLower(ask(r, "Proceed with setup? (y/n)", "y"))
			if confirm != "y" && confirm != "yes" {
				fmt.Println("Setup cancelled.")
				return nil
			}

			os.Args = []string{
				os.Args[0], "setup",
				"--account", acc,
				"--git-alias", gitAlias,
				"--git-host", gitHost,
				"--git-port", fmt.Sprintf("%d", gitPort),
				"--gitlab-group", gitlabGroup,
				"--folder", folder,
				"--name", name,
				"--email", email,
				"--key", key,
			}
			if passphrase {
				os.Args = append(os.Args, "--passphrase")
			}
			if uploadKey == "y" || uploadKey == "yes" {
				os.Args = append(os.Args, "--upload-key")
			}
			if adminAlias != "" && adminHost != "" {
				os.Args = append(os.Args, "--admin-alias", adminAlias, "--admin-host", adminHost, "--admin-port", fmt.Sprintf("%d", adminPort))
			}
			return rootCmd.Execute()
		},
	}
	rootCmd.AddCommand(wizardCmd)
}

func ask(r *bufio.Reader, prompt, def string) string {
	fmt.Printf("%s [%s]: ", prompt, def)
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" { return def }
	return line
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}