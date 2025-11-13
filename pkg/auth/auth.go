package auth

import (
	"fmt"
	"os"
	"path/filepath"
)

// TokenManager quản lý GitHub token
type TokenManager struct {
	token string
}

// NewTokenManager khởi tạo TokenManager
func NewTokenManager() *TokenManager {
	return &TokenManager{}
}

// GetToken lấy token từ environment hoặc file config
func (tm *TokenManager) GetToken(token string) (string, error) {
	if token != "" {
		tm.token = token
		return token, nil
	}

	// Kiểm tra environment variable
	if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
		tm.token = envToken
		return envToken, nil
	}

	// Kiểm tra file config
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "bug-crawler", "token")
	if data, err := os.ReadFile(configPath); err == nil {
		token := string(data)
		tm.token = token
		return token, nil
	}

	return "", fmt.Errorf("token không tìm thấy")
}

// SaveToken lưu token vào file config
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
