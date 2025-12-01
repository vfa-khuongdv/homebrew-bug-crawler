package auth

import (
	"fmt"
	"os"
	"path/filepath"
)

// TokenManager manages tokens for multiple platforms
type TokenManager struct {
	configDir string
}

// NewTokenManager creates a new TokenManager
func NewTokenManager() *TokenManager {
	return &TokenManager{
		configDir: filepath.Join(os.Getenv("HOME"), ".config", "bug-crawler"),
	}
}

// GetTokenForPlatform gets token for a specific platform
func (tm *TokenManager) GetTokenForPlatform(platform string) (string, error) {
	// Check environment variable (legacy support for GitHub)
	if platform == "github" {
		if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
			return envToken, nil
		}
	}

	// Check platform-specific config file
	tokenFile := filepath.Join(tm.configDir, platform+"_token")
	if data, err := os.ReadFile(tokenFile); err == nil {
		return string(data), nil
	}

	// Fallback: check legacy token file for GitHub
	if platform == "github" {
		legacyFile := filepath.Join(tm.configDir, "token")
		if data, err := os.ReadFile(legacyFile); err == nil {
			// Migrate to new format
			_ = tm.SaveTokenForPlatform(platform, string(data))
			return string(data), nil
		}
	}

	return "", fmt.Errorf("token không tìm thấy cho %s", platform)
}

// SaveTokenForPlatform saves token for a specific platform
func (tm *TokenManager) SaveTokenForPlatform(platform, token string) error {
	if err := os.MkdirAll(tm.configDir, 0700); err != nil {
		return err
	}

	tokenFile := filepath.Join(tm.configDir, platform+"_token")
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

// GetBacklogSpaceID gets Backlog space ID
func (tm *TokenManager) GetBacklogSpaceID() (string, error) {
	spaceFile := filepath.Join(tm.configDir, "backlog_space")
	if data, err := os.ReadFile(spaceFile); err == nil {
		return string(data), nil
	}
	return "", fmt.Errorf("backlog space ID không tìm thấy")
}

// SaveBacklogSpaceID saves Backlog space ID
func (tm *TokenManager) SaveBacklogSpaceID(spaceID string) error {
	if err := os.MkdirAll(tm.configDir, 0700); err != nil {
		return err
	}

	spaceFile := filepath.Join(tm.configDir, "backlog_space")
	return os.WriteFile(spaceFile, []byte(spaceID), 0600)
}

// GetBitbucketUsername gets Bitbucket username
func (tm *TokenManager) GetBitbucketUsername() (string, error) {
	usernameFile := filepath.Join(tm.configDir, "bitbucket_username")
	if data, err := os.ReadFile(usernameFile); err == nil {
		return string(data), nil
	}
	return "", fmt.Errorf("bitbucket username không tìm thấy")
}

// SaveBitbucketUsername saves Bitbucket username
func (tm *TokenManager) SaveBitbucketUsername(username string) error {
	if err := os.MkdirAll(tm.configDir, 0700); err != nil {
		return err
	}

	usernameFile := filepath.Join(tm.configDir, "bitbucket_username")
	return os.WriteFile(usernameFile, []byte(username), 0600)
}

// GetBacklogDomain gets Backlog domain
func (tm *TokenManager) GetBacklogDomain() (string, error) {
	domainFile := filepath.Join(tm.configDir, "backlog_domain")
	if data, err := os.ReadFile(domainFile); err == nil {
		return string(data), nil
	}
	return "", fmt.Errorf("backlog domain không tìm thấy")
}

// SaveBacklogDomain saves Backlog domain
func (tm *TokenManager) SaveBacklogDomain(domain string) error {
	if err := os.MkdirAll(tm.configDir, 0700); err != nil {
		return err
	}

	domainFile := filepath.Join(tm.configDir, "backlog_domain")
	return os.WriteFile(domainFile, []byte(domain), 0600)
}
