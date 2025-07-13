package vibelogger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMultiProjectLogOrganization(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	// Test default project (no project name specified)
	t.Run("DefaultProject", func(t *testing.T) {
		config := &LoggerConfig{
			AutoSave:        true,
			EnableMemoryLog: false,
			ProjectName:     "", // No project name specified
		}

		logger, err := CreateFileLoggerWithConfig("app", config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		// Check if file is created in logs/default/ directory
		expectedDir := filepath.Join("logs", "default")
		if !strings.HasPrefix(logger.filePath, expectedDir) {
			t.Errorf("Expected file path to start with %s, got %s", expectedDir, logger.filePath)
		}

		// Verify directory exists
		if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to exist", expectedDir)
		}

		// Write a test log
		err = logger.Info("test", "Default project log message")
		if err != nil {
			t.Fatalf("Failed to write log: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(logger.filePath); os.IsNotExist(err) {
			t.Errorf("Expected log file %s to exist", logger.filePath)
		}
	})

	// Test specific project name
	t.Run("SpecificProject", func(t *testing.T) {
		config := &LoggerConfig{
			AutoSave:        true,
			EnableMemoryLog: false,
			ProjectName:     "my-awesome-project",
		}

		logger, err := CreateFileLoggerWithConfig("service", config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		// Check if file is created in logs/my-awesome-project/ directory
		expectedDir := filepath.Join("logs", "my-awesome-project")
		if !strings.HasPrefix(logger.filePath, expectedDir) {
			t.Errorf("Expected file path to start with %s, got %s", expectedDir, logger.filePath)
		}

		// Verify directory exists
		if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to exist", expectedDir)
		}

		// Write a test log
		err = logger.Info("test", "Specific project log message")
		if err != nil {
			t.Fatalf("Failed to write log: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(logger.filePath); os.IsNotExist(err) {
			t.Errorf("Expected log file %s to exist", logger.filePath)
		}
	})

	// Test multiple projects
	t.Run("MultipleProjects", func(t *testing.T) {
		projects := []string{"project-a", "project-b", "project-c"}
		loggers := make([]*Logger, len(projects))

		// Create loggers for different projects
		for i, project := range projects {
			config := &LoggerConfig{
				AutoSave:        true,
				EnableMemoryLog: false,
				ProjectName:     project,
			}

			logger, err := CreateFileLoggerWithConfig("worker", config)
			if err != nil {
				t.Fatalf("Failed to create logger for project %s: %v", project, err)
			}
			loggers[i] = logger

			// Write project-specific log
			err = logger.Info("startup", "Worker started for project "+project)
			if err != nil {
				t.Fatalf("Failed to write log for project %s: %v", project, err)
			}
		}

		// Verify all projects have separate directories
		for i, project := range projects {
			expectedDir := filepath.Join("logs", project)
			if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
				t.Errorf("Expected directory %s to exist", expectedDir)
			}

			// Verify file path is correct
			if !strings.HasPrefix(loggers[i].filePath, expectedDir) {
				t.Errorf("Expected file path to start with %s, got %s", expectedDir, loggers[i].filePath)
			}
		}

		// Clean up
		for _, logger := range loggers {
			logger.Close()
		}
	})
}

func TestProjectNameValidation(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		valid       bool
	}{
		{"ValidAlphanumeric", "project123", true},
		{"ValidWithUnderscore", "my_project", true},
		{"ValidWithHyphen", "my-project", true},
		{"ValidMixed", "Project_123-test", true},
		{"InvalidWithDot", "project.name", false},
		{"InvalidWithSpace", "my project", false},
		{"InvalidWithSlash", "project/name", false},
		{"InvalidWithSpecialChar", "project@name", false},
		{"EmptyString", "", true}, // Empty is valid (uses default)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidProjectName(tt.projectName)
			if result != tt.valid {
				t.Errorf("isValidProjectName(%s) = %v, want %v", tt.projectName, result, tt.valid)
			}
		})
	}
}

func TestProjectNameEnvironmentVariable(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
		os.Unsetenv("VIBE_LOG_PROJECT_NAME")
	}()

	// Test valid project name from environment
	t.Run("ValidEnvironmentProjectName", func(t *testing.T) {
		os.Setenv("VIBE_LOG_PROJECT_NAME", "env-project")

		config, err := NewConfigFromEnvironment()
		if err != nil {
			t.Fatalf("Failed to create config from environment: %v", err)
		}

		if config.ProjectName != "env-project" {
			t.Errorf("Expected project name 'env-project', got '%s'", config.ProjectName)
		}

		// Test logger creation with environment config
		logger, err := CreateFileLoggerWithConfig("app", config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		expectedDir := filepath.Join("logs", "env-project")
		if !strings.HasPrefix(logger.filePath, expectedDir) {
			t.Errorf("Expected file path to start with %s, got %s", expectedDir, logger.filePath)
		}
	})

	// Test invalid project name from environment
	t.Run("InvalidEnvironmentProjectName", func(t *testing.T) {
		os.Setenv("VIBE_LOG_PROJECT_NAME", "invalid@project")

		_, err := NewConfigFromEnvironment()
		if err == nil {
			t.Error("Expected error for invalid project name, got nil")
		}

		if !strings.Contains(err.Error(), "VIBE_LOG_PROJECT_NAME contains invalid characters") {
			t.Errorf("Expected error about invalid characters, got: %v", err)
		}
	})
}

func TestCustomFilePathOverridesProject(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
		os.RemoveAll("custom")
	}()

	config := &LoggerConfig{
		AutoSave:        true,
		EnableMemoryLog: false,
		ProjectName:     "should-be-ignored",
		FilePath:        "custom/specific.log",
	}

	logger, err := CreateFileLoggerWithConfig("app", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Custom file path should override project organization
	if logger.filePath != "custom/specific.log" {
		t.Errorf("Expected file path 'custom/specific.log', got '%s'", logger.filePath)
	}

	// Project directory should not be created
	projectDir := filepath.Join("logs", "should-be-ignored")
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		t.Errorf("Project directory %s should not exist when custom file path is used", projectDir)
	}

	// Custom directory should exist
	if _, err := os.Stat("custom"); os.IsNotExist(err) {
		t.Errorf("Custom directory should exist")
	}
}