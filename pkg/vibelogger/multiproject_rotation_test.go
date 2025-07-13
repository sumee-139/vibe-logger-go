package vibelogger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMultiProjectRotationIntegration(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	t.Run("ProjectRotationWithSpecificProject", func(t *testing.T) {
		config := &LoggerConfig{
			MaxFileSize:     100, // Small for testing
			AutoSave:        true,
			RotationEnabled: true,
			MaxRotatedFiles: 3,
			ProjectName:     "rotation-test-project",
		}

		if err := os.MkdirAll("test_logs", 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		config.FilePath = "test_logs/rotation-test.log"

		logger, err := CreateFileLoggerWithConfig("rotation-test", config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		// Write logs to trigger rotation
		for i := 0; i < 10; i++ {
			err := logger.Info("rotation_test", "This is a test message for rotation in project context")
			if err != nil {
				t.Fatalf("Failed to write log entry %d: %v", i, err)
			}
		}

		// Force rotation to ensure it happens
		err = logger.ForceRotation()
		if err != nil {
			t.Fatalf("Failed to force rotation: %v", err)
		}

		// Check if rotated files exist
		rotatedFiles := logger.GetRotatedFiles()
		if len(rotatedFiles) == 0 {
			t.Error("Expected at least one rotated file")
		}

		// Verify all rotated files are in the correct directory structure
		for _, file := range rotatedFiles {
			if !strings.Contains(file, "test_logs") {
				t.Errorf("Rotated file %s should be in test_logs directory", file)
			}
		}
	})

	t.Run("ProjectRotationWithDefaultStructure", func(t *testing.T) {
		config := &LoggerConfig{
			MaxFileSize:     100, // Small for testing
			AutoSave:        true,
			RotationEnabled: true,
			MaxRotatedFiles: 2,
			ProjectName:     "auto-project", // Should create logs/auto-project/
		}

		logger, err := CreateFileLoggerWithConfig("auto-service", config)
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		defer logger.Close()

		// Verify the logger was created in the correct project directory
		expectedDir := filepath.Join("logs", "auto-project")
		if !strings.HasPrefix(logger.filePath, expectedDir) {
			t.Errorf("Expected file path to start with %s, got %s", expectedDir, logger.filePath)
		}

		// Write logs to trigger rotation
		for i := 0; i < 5; i++ {
			err := logger.Info("service_test", "Auto-rotation test message in project directory")
			if err != nil {
				t.Fatalf("Failed to write log entry %d: %v", i, err)
			}
		}

		// Force rotation
		err = logger.ForceRotation()
		if err != nil {
			t.Fatalf("Failed to force rotation: %v", err)
		}

		// Check rotated files
		rotatedFiles := logger.GetRotatedFiles()
		if len(rotatedFiles) == 0 {
			t.Error("Expected at least one rotated file")
		}

		// Verify all rotated files are in the project directory
		for _, file := range rotatedFiles {
			if !strings.HasPrefix(file, expectedDir) {
				t.Errorf("Rotated file %s should be in project directory %s", file, expectedDir)
			}
		}

		// Verify project directory exists and contains files
		if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
			t.Errorf("Project directory %s should exist", expectedDir)
		}
	})

	t.Run("MultipleProjectsWithIndependentRotation", func(t *testing.T) {
		projects := []string{"project-alpha", "project-beta"}
		loggers := make([]*Logger, len(projects))

		// Create loggers for different projects
		for i, project := range projects {
			config := &LoggerConfig{
				MaxFileSize:     80, // Small for testing
				AutoSave:        true,
				RotationEnabled: true,
				MaxRotatedFiles: 2,
				ProjectName:     project,
			}

			logger, err := CreateFileLoggerWithConfig("worker", config)
			if err != nil {
				t.Fatalf("Failed to create logger for project %s: %v", project, err)
			}
			loggers[i] = logger

			// Write logs for each project
			for j := 0; j < 5; j++ {
				err = logger.Info("worker_test", "Worker message for "+project)
				if err != nil {
					t.Fatalf("Failed to write log for project %s: %v", project, err)
				}
			}

			// Force rotation for each project
			err = logger.ForceRotation()
			if err != nil {
				t.Fatalf("Failed to force rotation for project %s: %v", project, err)
			}
		}

		// Verify each project has its own rotation files
		for i, project := range projects {
			logger := loggers[i]
			rotatedFiles := logger.GetRotatedFiles()

			if len(rotatedFiles) == 0 {
				t.Errorf("Project %s should have rotated files", project)
			}

			expectedDir := filepath.Join("logs", project)
			for _, file := range rotatedFiles {
				if !strings.HasPrefix(file, expectedDir) {
					t.Errorf("Rotated file %s for project %s should be in directory %s", file, project, expectedDir)
				}

				// Verify files don't interfere with other projects
				for _, otherProject := range projects {
					if otherProject != project {
						otherDir := filepath.Join("logs", otherProject)
						if strings.HasPrefix(file, otherDir) {
							t.Errorf("Rotated file %s for project %s should not be in other project directory %s", file, project, otherDir)
						}
					}
				}
			}
		}

		// Clean up
		for _, logger := range loggers {
			logger.Close()
		}
	})
}

func TestProjectDefaultDirectoryFallback(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
		os.RemoveAll("logs")
	}()

	// Test that empty project name creates default directory
	config := &LoggerConfig{
		AutoSave:        true,
		RotationEnabled: true,
		MaxRotatedFiles: 2,
		ProjectName:     "", // Empty project name
	}

	logger, err := CreateFileLoggerWithConfig("default-test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Verify file is in logs/default/ directory
	expectedDir := filepath.Join("logs", "default")
	if !strings.HasPrefix(logger.filePath, expectedDir) {
		t.Errorf("Expected file path to start with %s, got %s", expectedDir, logger.filePath)
	}

	// Verify directory exists
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Default directory %s should exist", expectedDir)
	}

	// Test rotation works in default directory
	for i := 0; i < 3; i++ {
		err := logger.Info("default_test", "Message in default project directory")
		if err != nil {
			t.Fatalf("Failed to write log: %v", err)
		}
	}

	// Force rotation
	err = logger.ForceRotation()
	if err != nil {
		t.Fatalf("Failed to force rotation: %v", err)
	}

	// Verify rotated files are in default directory
	rotatedFiles := logger.GetRotatedFiles()
	for _, file := range rotatedFiles {
		if !strings.HasPrefix(file, expectedDir) {
			t.Errorf("Rotated file %s should be in default directory %s", file, expectedDir)
		}
	}
}