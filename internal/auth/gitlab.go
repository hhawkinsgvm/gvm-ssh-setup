package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

// GitLabConfig holds GitLab connection configuration
type GitLabConfig struct {
	BaseURL string
	Token   string
	Client  *gitlab.Client
}

// Group represents a GitLab group
type Group struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}

// NewGitLabConfig creates a new GitLab configuration
func NewGitLabConfig(baseURL, token string) (*GitLabConfig, error) {
	if token == "" {
		return nil, fmt.Errorf("GitLab token is required")
	}

	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}

	// Create GitLab client
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	return &GitLabConfig{
		BaseURL: baseURL,
		Token:   token,
		Client:  client,
	}, nil
}

// CheckGroupMembership verifies if the user is a member of the specified group
func (gc *GitLabConfig) CheckGroupMembership(groupPath string) (bool, error) {
	// Get current user
	user, _, err := gc.Client.Users.CurrentUser()
	if err != nil {
		return false, fmt.Errorf("failed to get current user: %w", err)
	}

	// Get user's groups
	groups, _, err := gc.Client.Groups.ListGroups(&gitlab.ListGroupsOptions{
		MinAccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
	})
	if err != nil {
		return false, fmt.Errorf("failed to list user groups: %w", err)
	}

	// Check if user is member of the target group
	for _, group := range groups {
		if group.FullPath == groupPath || group.Path == groupPath {
			return true, nil
		}
	}

	// Also check if user is directly a member of the group
	_, _, err = gc.Client.GroupMembers.GetGroupMember(groupPath, user.ID)
	if err == nil {
		return true, nil
	}

	return false, nil
}

// AddSSHKey adds an SSH key to the user's GitLab account
func (gc *GitLabConfig) AddSSHKey(title, keyContent string) error {
	// Check if key already exists
	keys, _, err := gc.Client.Users.ListSSHKeys(&gitlab.ListSSHKeysOptions{})
	if err != nil {
		return fmt.Errorf("failed to list existing SSH keys: %w", err)
	}

	// Extract the public key part (remove key type and comment)
	keyParts := strings.Fields(keyContent)
	if len(keyParts) < 2 {
		return fmt.Errorf("invalid SSH key format")
	}
	keyData := keyParts[1] // The actual key data without ssh-ed25519 prefix

	for _, key := range keys {
		if strings.Contains(key.Key, keyData) {
			return fmt.Errorf("SSH key already exists with title: %s", key.Title)
		}
	}

	// Add the new key
	_, _, err = gc.Client.Users.AddSSHKey(&gitlab.AddSSHKeyOptions{
		Title: &title,
		Key:   &keyContent,
	})
	if err != nil {
		return fmt.Errorf("failed to add SSH key: %w", err)
	}

	return nil
}

// GetGitLabToken attempts to get GitLab token from environment or glab CLI
func GetGitLabToken() string {
	// First try environment variable
	if token := os.Getenv("GITLAB_TOKEN"); token != "" {
		return token
	}

	// Try to get from glab CLI config
	if token := getGlabToken(); token != "" {
		return token
	}

	return ""
}

// getGlabToken attempts to extract token from glab CLI configuration
func getGlabToken() string {
	// This is a simplified approach - in reality, glab stores tokens
	// in a more complex format. For now, we'll rely on GITLAB_TOKEN env var.
	return ""
}

// ValidateGitLabAccess validates that we can access GitLab with the given configuration
func ValidateGitLabAccess(baseURL, token string) error {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", baseURL+"/api/v4/user", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to GitLab: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("GitLab API returned status %d", resp.StatusCode)
	}

	return nil
}
