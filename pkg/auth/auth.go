package auth

import (
	"fmt"
	"os"
	"path/filepath"
)

// TokenManager manages GitHub token
type TokenManager struct {
	token string
}

// NewTokenManager creates a new TokenManager
func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

// GetToken gets token from environment or config file
func (tm *TokenManager) GetToken(token string) (string, error) {
	if token != "" {
		tm.token = token
		return token, nil
	}

	// Check environment variable
	if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
		tm.token = envToken
		return envToken, nil
	}

	// Check config file
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "bug-crawler", "token")
	if data, err := os.ReadFile(configPath); err == nil {
		token := string(data)
		tm.token = token
		return token, nil
	}

	return "", fmt.Errorf("token không tìm thấy")
}

// SaveToken saves token to config file
func (tm *TokenManager) SaveToken(token string) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "bug-crawler")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "token")
	return os.WriteFile(configFile, []byte(token), 0600)
}

// Token getter
func (tm *TokenManager) Token() string {
	return tm.token
}
