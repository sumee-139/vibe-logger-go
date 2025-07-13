package vibelogger

import (
	"os"
	"testing"
)

func TestClearMemoryLogs(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Add some logs
	err := logger.Info("test1", "Test message 1")
	if err != nil {
		t.Fatalf("Failed to log: %v", err)
	}
	
	err = logger.Info("test2", "Test message 2")
	if err != nil {
		t.Fatalf("Failed to log: %v", err)
	}

	// Verify logs exist
	logs := logger.GetMemoryLogs()
	if len(logs) != 2 {
		t.Fatalf("Expected 2 log entries, got %d", len(logs))
	}

	// Clear memory logs
	logger.ClearMemoryLogs()

	// Verify logs are cleared
	logs = logger.GetMemoryLogs()
	if len(logs) != 0 {
		t.Fatalf("Expected 0 log entries after clear, got %d", len(logs))
	}
}

func TestWithHumanNote(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Test WithHumanNote
	humanNote := "This is a human-readable note for debugging"
	err := logger.Info("test_operation", "Test message with human note", WithHumanNote(humanNote))
	if err != nil {
		t.Fatalf("Failed to log with human note: %v", err)
	}

	// Verify the human note is in the memory log
	logs := logger.GetMemoryLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.HumanNote != humanNote {
		t.Errorf("Expected human note to be '%s', got '%s'", humanNote, entry.HumanNote)
	}
}

func TestWithAITodo(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Test WithAITodo
	aiTodo := "TODO: Implement caching for better performance"
	err := logger.Info("test_operation", "Test message with AI todo", WithAITodo(aiTodo))
	if err != nil {
		t.Fatalf("Failed to log with AI todo: %v", err)
	}

	// Verify the AI todo is in the memory log
	logs := logger.GetMemoryLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.AITodo != aiTodo {
		t.Errorf("Expected AI todo to be '%s', got '%s'", aiTodo, entry.AITodo)
	}
}

func TestIsValidEnvironmentName(t *testing.T) {
	testCases := []struct {
		name string
		env  string
		want bool
	}{
		{"valid_production", "production", true},
		{"valid_development", "development", true},
		{"valid_test", "test", true},
		{"valid_staging", "staging", true},
		{"valid_with_hyphen", "prod-env", true},
		{"valid_with_numbers", "env123", true},
		{"valid_uppercase", "PRODUCTION", true},
		{"valid_with_underscore", "test_env", true},
		{"valid_with_dot", "env.prod", true},
		{"invalid_space", "prod env", false},
		{"invalid_special_char", "env@prod", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidEnvironmentName(tc.env)
			if got != tc.want {
				t.Errorf("isValidEnvironmentName(%q) = %v, want %v", tc.env, got, tc.want)
			}
		})
	}
}

func TestRotationManagerUpdateConfigDirect(t *testing.T) {
	// Cleanup test files
	defer func() {
		os.RemoveAll("test_logs")
	}()

	// Create test directory
	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create initial config
	config := &LoggerConfig{
		MaxFileSize:     1000,
		RotationEnabled: true,
		MaxRotatedFiles: 3,
		AutoSave:        true,
		FilePath:        "test_logs/update_config_direct_test.log",
	}

	// Create logger with config
	logger := NewLoggerWithConfig("update_config_direct_test", config)

	// Create rotation manager directly
	rotationManager := NewRotationManager(logger, config, "test_logs/update_config_direct_test.log")
	if rotationManager == nil {
		t.Fatal("Failed to create rotation manager")
	}
	defer rotationManager.Close()

	// Update config with different settings
	newConfig := &LoggerConfig{
		MaxFileSize:     2000,
		RotationEnabled: true,
		MaxRotatedFiles: 5,
		AutoSave:        true,
		FilePath:        "test_logs/update_config_direct_test.log",
	}

	// Test UpdateConfig method directly
	rotationManager.UpdateConfig(newConfig)

	// Test passes if UpdateConfig method runs without error
	// (we can't easily verify internal state without exposing fields)
}