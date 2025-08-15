package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
)

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	config := DefaultConfig()

	// Load from environment variables
	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}

	if username := os.Getenv("ADMIN_USERNAME"); username != "" {
		config.AdminUsername = username
	}

	if password := os.Getenv("ADMIN_PASSWORD"); password != "" {
		// Hash the password if it's provided as plain text
		config.AdminPassword = hashPassword(password)
	}

	if maxSizeStr := os.Getenv("UPLOAD_MAX_SIZE"); maxSizeStr != "" {
		if maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64); err == nil {
			config.UploadMaxSize = maxSize
		}
	}

	if timeoutStr := os.Getenv("SESSION_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			config.SessionTimeout = timeout
		}
	}

	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		config.DataDir = dataDir
	}

	if staticDir := os.Getenv("STATIC_DIR"); staticDir != "" {
		config.StaticDir = staticDir
	}

	if templatesDir := os.Getenv("TEMPLATES_DIR"); templatesDir != "" {
		config.TemplatesDir = templatesDir
	}

	return config
}

// ValidateConfig validates the configuration
func ValidateConfig(config *Config) error {
	// Basic validation - can be expanded later
	if config.Port == "" {
		config.Port = "8080"
	}

	if config.AdminUsername == "" {
		config.AdminUsername = "admin"
	}

	if config.AdminPassword == "" {
		// Hash the default password
		config.AdminPassword = hashPassword("admin123")
	}

	return nil
}

// hashPassword creates a SHA-256 hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
