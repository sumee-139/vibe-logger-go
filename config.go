package vibelogger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Security and resource limits
const (
	MaxFileSizeLimit  = 1 * 1024 * 1024 * 1024 // 1GB maximum file size
	MaxMemoryLogLimit = 10000                  // 10k entries maximum
	MaxFilePathLength = 255                    // 255 characters maximum
)

// LoggerConfig represents configuration options for the logger
type LoggerConfig struct {
	MaxFileSize     int64  `json:"max_file_size"`     // Maximum file size in bytes (0 = unlimited)
	AutoSave        bool   `json:"auto_save"`         // Enable/disable auto-save functionality
	EnableMemoryLog bool   `json:"enable_memory_log"` // Enable in-memory logging
	MemoryLogLimit  int    `json:"memory_log_limit"`  // Maximum number of entries in memory log
	FilePath        string `json:"file_path"`         // Custom log file path
	Environment     string `json:"environment"`       // Environment name (dev/prod/test)
	ProjectName     string `json:"project_name"`      // Project name for multi-project log organization
	// Log rotation settings
	RotationEnabled bool `json:"rotation_enabled"`  // Enable/disable log rotation
	MaxRotatedFiles int  `json:"max_rotated_files"` // Maximum number of rotated files to keep (0 = keep all)
}

// DefaultConfig returns a LoggerConfig with sensible defaults
func DefaultConfig() *LoggerConfig {
	return &LoggerConfig{
		MaxFileSize:     10 * 1024 * 1024, // 10MB default
		AutoSave:        true,             // Auto-save enabled by default
		EnableMemoryLog: false,            // Memory log disabled by default
		MemoryLogLimit:  1000,             // 1000 entries default
		FilePath:        "",               // Use default path generation
		Environment:     "development",    // Default environment
		ProjectName:     "",               // Use default project organization
		RotationEnabled: true,             // Log rotation enabled by default
		MaxRotatedFiles: 5,                // Keep 5 rotated files by default
	}
}

// LoadFromEnvironment loads configuration from environment variables with validation
func (c *LoggerConfig) LoadFromEnvironment() error {
	var validationErrors []string

	// Validate VIBE_LOG_MAX_FILE_SIZE
	if val := os.Getenv("VIBE_LOG_MAX_FILE_SIZE"); val != "" {
		if size, err := strconv.ParseInt(val, 10, 64); err == nil {
			if size < 0 {
				validationErrors = append(validationErrors, "VIBE_LOG_MAX_FILE_SIZE cannot be negative")
			} else if size > MaxFileSizeLimit {
				validationErrors = append(validationErrors, fmt.Sprintf("VIBE_LOG_MAX_FILE_SIZE exceeds limit: %d > %d", size, MaxFileSizeLimit))
			} else {
				c.MaxFileSize = size
			}
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid VIBE_LOG_MAX_FILE_SIZE format: %s", val))
		}
	}

	// Validate VIBE_LOG_AUTO_SAVE
	if val := os.Getenv("VIBE_LOG_AUTO_SAVE"); val != "" {
		if autoSave, err := strconv.ParseBool(val); err == nil {
			c.AutoSave = autoSave
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid VIBE_LOG_AUTO_SAVE format: %s (must be true/false)", val))
		}
	}

	// Validate VIBE_LOG_ENABLE_MEMORY
	if val := os.Getenv("VIBE_LOG_ENABLE_MEMORY"); val != "" {
		if enableMemory, err := strconv.ParseBool(val); err == nil {
			c.EnableMemoryLog = enableMemory
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid VIBE_LOG_ENABLE_MEMORY format: %s (must be true/false)", val))
		}
	}

	// Validate VIBE_LOG_MEMORY_LIMIT
	if val := os.Getenv("VIBE_LOG_MEMORY_LIMIT"); val != "" {
		if limit, err := strconv.Atoi(val); err == nil {
			if limit < 0 {
				validationErrors = append(validationErrors, "VIBE_LOG_MEMORY_LIMIT cannot be negative")
			} else if limit > MaxMemoryLogLimit {
				validationErrors = append(validationErrors, fmt.Sprintf("VIBE_LOG_MEMORY_LIMIT exceeds limit: %d > %d", limit, MaxMemoryLogLimit))
			} else {
				c.MemoryLogLimit = limit
			}
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid VIBE_LOG_MEMORY_LIMIT format: %s", val))
		}
	}

	// Validate VIBE_LOG_FILE_PATH
	if val := os.Getenv("VIBE_LOG_FILE_PATH"); val != "" {
		if len(val) > MaxFilePathLength {
			validationErrors = append(validationErrors, fmt.Sprintf("VIBE_LOG_FILE_PATH too long: %d > %d", len(val), MaxFilePathLength))
		} else {
			// Temporarily set to validate path security
			oldPath := c.FilePath
			c.FilePath = val
			if err := c.validateFilePath(); err != nil {
				validationErrors = append(validationErrors, fmt.Sprintf("VIBE_LOG_FILE_PATH validation failed: %v", err))
				c.FilePath = oldPath // Restore old path on error
			}
		}
	}

	// Validate VIBE_LOG_ENVIRONMENT
	if val := os.Getenv("VIBE_LOG_ENVIRONMENT"); val != "" {
		// Environment names should be reasonable length and safe characters
		if len(val) > 50 {
			validationErrors = append(validationErrors, "VIBE_LOG_ENVIRONMENT too long (max 50 characters)")
		} else if !isValidEnvironmentName(val) {
			validationErrors = append(validationErrors, fmt.Sprintf("VIBE_LOG_ENVIRONMENT contains invalid characters: %s", val))
		} else {
			c.Environment = val
		}
	}

	// Validate VIBE_LOG_PROJECT_NAME
	if val := os.Getenv("VIBE_LOG_PROJECT_NAME"); val != "" {
		// Project names should be reasonable length and safe characters
		if len(val) > 50 {
			validationErrors = append(validationErrors, "VIBE_LOG_PROJECT_NAME too long (max 50 characters)")
		} else if !isValidProjectName(val) {
			validationErrors = append(validationErrors, fmt.Sprintf("VIBE_LOG_PROJECT_NAME contains invalid characters: %s", val))
		} else {
			c.ProjectName = val
		}
	}

	// Validate VIBE_LOG_ROTATION_ENABLED
	if val := os.Getenv("VIBE_LOG_ROTATION_ENABLED"); val != "" {
		if rotation, err := strconv.ParseBool(val); err == nil {
			c.RotationEnabled = rotation
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid VIBE_LOG_ROTATION_ENABLED format: %s (must be true/false)", val))
		}
	}

	// Validate VIBE_LOG_MAX_ROTATED_FILES
	if val := os.Getenv("VIBE_LOG_MAX_ROTATED_FILES"); val != "" {
		if files, err := strconv.Atoi(val); err == nil {
			if files < 0 {
				validationErrors = append(validationErrors, "VIBE_LOG_MAX_ROTATED_FILES cannot be negative")
			} else if files > 100 {
				validationErrors = append(validationErrors, "VIBE_LOG_MAX_ROTATED_FILES too large (max 100)")
			} else {
				c.MaxRotatedFiles = files
			}
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("invalid VIBE_LOG_MAX_ROTATED_FILES format: %s", val))
		}
	}

	// Return validation errors if any
	if len(validationErrors) > 0 {
		return fmt.Errorf("environment variable validation errors: %v", validationErrors)
	}

	return nil
}

// isValidEnvironmentName checks if environment name contains only safe characters
func isValidEnvironmentName(env string) bool {
	// Allow alphanumeric, underscore, hyphen, and dot
	for _, char := range env {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-' || char == '.') {
			return false
		}
	}
	return true
}

// isValidProjectName checks if project name contains only safe characters
func isValidProjectName(project string) bool {
	// Allow alphanumeric, underscore, and hyphen (no dots for directory safety)
	for _, char := range project {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}
	return true
}

// NewConfigFromEnvironment creates a new LoggerConfig with environment variables applied
func NewConfigFromEnvironment() (*LoggerConfig, error) {
	config := DefaultConfig()
	if err := config.LoadFromEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to load configuration from environment: %w", err)
	}
	return config, nil
}

// Validate checks if the configuration is valid and secure
func (c *LoggerConfig) Validate() error {
	// Validate file size limits
	if c.MaxFileSize < 0 {
		c.MaxFileSize = 0 // 0 means unlimited
	}
	if c.MaxFileSize > MaxFileSizeLimit {
		return fmt.Errorf("max file size exceeds limit: %d > %d", c.MaxFileSize, MaxFileSizeLimit)
	}

	// Validate memory log limits
	if c.MemoryLogLimit < 0 {
		c.MemoryLogLimit = 0 // 0 means unlimited
	}
	if c.MemoryLogLimit > MaxMemoryLogLimit {
		return fmt.Errorf("memory log limit exceeds maximum: %d > %d", c.MemoryLogLimit, MaxMemoryLogLimit)
	}

	// Validate file path security
	if err := c.validateFilePath(); err != nil {
		return fmt.Errorf("file path validation failed: %w", err)
	}

	// Set default environment if empty
	if c.Environment == "" {
		c.Environment = "development"
	}

	return nil
}

// validateFilePath ensures the file path is secure and prevents path traversal attacks
func (c *LoggerConfig) validateFilePath() error {
	if c.FilePath == "" {
		return nil // Empty path is okay, will use default
	}

	// Check path length
	if len(c.FilePath) > MaxFilePathLength {
		return fmt.Errorf("file path too long: %d > %d characters", len(c.FilePath), MaxFilePathLength)
	}

	// Prevent path traversal attacks
	if strings.Contains(c.FilePath, "..") {
		return fmt.Errorf("file path contains path traversal characters: %s", c.FilePath)
	}

	// Clean the path to normalize it
	cleanPath := filepath.Clean(c.FilePath)

	// For relative paths, ensure they're within safe directories
	if !filepath.IsAbs(cleanPath) {
		// Only allow relative paths within ./logs/ directory or current directory
		if !strings.HasPrefix(cleanPath, "logs/") && !strings.HasPrefix(cleanPath, "./logs/") && cleanPath != "." {
			return fmt.Errorf("relative file path must be within logs directory: %s", cleanPath)
		}
	} else {
		// For absolute paths, only allow specific safe directories
		safeDirs := []string{"/tmp/", "/var/log/", "/home/"}
		allowed := false
		for _, safeDir := range safeDirs {
			if strings.HasPrefix(cleanPath, safeDir) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("absolute file path not in allowed directories (/tmp/, /var/log/, /home/): %s", cleanPath)
		}
	}

	// Update the path to the cleaned version
	c.FilePath = cleanPath

	return nil
}
