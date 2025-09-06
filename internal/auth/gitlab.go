package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// GitLabGroup represents a GitLab group from the API
type GitLabGroup struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}

// GitLabUser represents a GitLab user from the API
type GitLabUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// GitLabGroupMember checks if the current user is a member of the specified group
func GitLabGroupMember(apiBase, targetGroup string) (bool, error) {
	tok := os.Getenv("GITLAB_TOKEN")
	if tok == "" {
		return false, fmt.Errorf("GITLAB_TOKEN not set; run `glab auth login` or export a token")
	}
	
	// First, get user info to verify token is valid
	req, err := http.NewRequest("GET", apiBase+"/api/v4/user", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("PRIVATE-TOKEN", tok)
	
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("gitlab user check failed: %s", resp.Status)
	}
	
	var user GitLabUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return false, err
	}
	
	// Now check group membership
	req, err = http.NewRequest("GET", apiBase+"/api/v4/groups?min_access_level=10", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("PRIVATE-TOKEN", tok)
	
	resp, err = c.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("gitlab group check failed: %s", resp.Status)
	}
	
	var groups []GitLabGroup
	if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
		return false, err
	}
	
	for _, g := range groups {
		if g.FullPath == targetGroup || g.Path == targetGroup {
			return true, nil
		}
	}
	
	return false, nil
}

// CheckGitLabAccess verifies the user has access to the specified GitLab instance and group
func CheckGitLabAccess(gitHost, targetGroup string) error {
	// Try to determine API base from git host
	apiBase := "https://" + gitHost
	
	ok, err := GitLabGroupMember(apiBase, targetGroup)
	if err != nil {
		return fmt.Errorf("GitLab access check failed: %w", err)
	}
	
	if !ok {
		return fmt.Errorf("access denied: you are not a member of group '%s'", targetGroup)
	}
	
	return nil
}