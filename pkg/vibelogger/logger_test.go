package vibelogger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test_app")

	if logger.name != "test_app" {
		t.Errorf("Expected name 'test_app', got '%s'", logger.name)
	}

	if logger.file != nil {
		t.Error("Expected file to be nil for memory logger")
	}
}

func TestCreateFileLogger(t *testing.T) {
	// Clean up any existing logs directory
	defer func() {
		os.RemoveAll("logs")
	}()

	logger, err := CreateFileLogger("test_app")
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	defer logger.Close()

	if logger.name != "test_app" {
		t.Errorf("Expected name 'test_app', got '%s'", logger.name)
	}

	if logger.file == nil {
		t.Error("Expected file to be set for file logger")
	}

	// Check if log file exists
	if _, err := os.Stat(logger.filePath); os.IsNotExist(err) {
		t.Error("Log file should exist")
	}

	// Check if logs directory exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		t.Error("Logs directory should exist")
	}
}

func TestLogLevels(t *testing.T) {
	// Clean up any existing logs directory
	defer func() {
		os.RemoveAll("logs")
	}()

	logger, err := CreateFileLogger("test_levels")
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	defer logger.Close()

	tests := []struct {
		name      string
		logFunc   func(string, string, ...LogOption) error
		level     LogLevel
		operation string
		message   string
	}{
		{"Info", logger.Info, INFO, "test_op", "test info message"},
		{"Warn", logger.Warn, WARN, "test_op", "test warn message"},
		{"Error", logger.Error, ERROR, "test_op", "test error message"},
		{"Debug", logger.Debug, DEBUG, "test_op", "test debug message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.logFunc(tt.operation, tt.message)
			if err != nil {
				t.Errorf("Failed to log %s: %v", tt.name, err)
			}
		})
	}
}

func TestLogWithOptions(t *testing.T) {
	// Clean up any existing logs directory
	defer func() {
		os.RemoveAll("logs")
	}()

	logger, err := CreateFileLogger("test_options")
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	defer logger.Close()

	// Test with context
	err = logger.Info("test_context", "test message",
		WithContext(map[string]interface{}{
			"user_id": "123",
			"action":  "login",
		}),
	)
	if err != nil {
		t.Errorf("Failed to log with context: %v", err)
	}

	// Test with human note
	err = logger.Info("test_human_note", "test message",
		WithHumanNote("This is a test note"),
	)
	if err != nil {
		t.Errorf("Failed to log with human note: %v", err)
	}

	// Test with AI todo
	err = logger.Info("test_ai_todo", "test message",
		WithAITodo("AI-TODO: Test this functionality"),
	)
	if err != nil {
		t.Errorf("Failed to log with AI todo: %v", err)
	}

	// Test with correlation ID
	err = logger.Info("test_correlation", "test message",
		WithCorrelationID("test-123"),
	)
	if err != nil {
		t.Errorf("Failed to log with correlation ID: %v", err)
	}

	// Test with duration
	err = logger.Info("test_duration", "test message",
		WithDuration(100*time.Millisecond),
	)
	if err != nil {
		t.Errorf("Failed to log with duration: %v", err)
	}
}

func TestLoggerClose(t *testing.T) {
	// Clean up any existing logs directory
	defer func() {
		os.RemoveAll("logs")
	}()

	logger, err := CreateFileLogger("test_close")
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	// Write a log entry
	err = logger.Info("test_op", "test message")
	if err != nil {
		t.Errorf("Failed to log message: %v", err)
	}

	// Close the logger
	err = logger.Close()
	if err != nil {
		t.Errorf("Failed to close logger: %v", err)
	}

	// Try to close again (should not error)
	err = logger.Close()
	if err != nil {
		t.Errorf("Failed to close logger twice: %v", err)
	}
}

func TestGetEnvironment(t *testing.T) {
	env := getEnvironment()

	expectedKeys := []string{"go_version", "os", "arch", "pid", "pwd"}
	for _, key := range expectedKeys {
		if _, exists := env[key]; !exists {
			t.Errorf("Expected environment key '%s' to exist", key)
		}
	}
}

func TestGetStackTrace(t *testing.T) {
	stack := getStackTrace()

	if len(stack) == 0 {
		t.Error("Expected non-empty stack trace")
	}

	// Check that the first stack frame contains this test function
	found := false
	for _, frame := range stack {
		if filepath.Base(frame) != "" && len(frame) > 0 {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected stack trace to contain valid frame information")
	}
}

func TestLogEntryStructure(t *testing.T) {
	entry := LogEntry{
		Timestamp:     time.Now().UTC(),
		Level:         INFO,
		Operation:     "test_op",
		Message:       "test message",
		Context:       map[string]interface{}{"key": "value"},
		HumanNote:     "test note",
		AITodo:        "test todo",
		StackTrace:    []string{"frame1", "frame2"},
		Environment:   map[string]string{"key": "value"},
		CorrelationID: "test-123",
	}

	// Test that all fields are accessible
	if entry.Level != INFO {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}

	if entry.Operation != "test_op" {
		t.Errorf("Expected operation 'test_op', got '%s'", entry.Operation)
	}

	if entry.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", entry.Message)
	}

	if entry.Context["key"] != "value" {
		t.Errorf("Expected context key 'key' to have value 'value'")
	}

	if entry.HumanNote != "test note" {
		t.Errorf("Expected human note 'test note', got '%s'", entry.HumanNote)
	}

	if entry.AITodo != "test todo" {
		t.Errorf("Expected AI todo 'test todo', got '%s'", entry.AITodo)
	}

	if len(entry.StackTrace) != 2 {
		t.Errorf("Expected 2 stack frames, got %d", len(entry.StackTrace))
	}

	if entry.Environment["key"] != "value" {
		t.Errorf("Expected environment key 'key' to have value 'value'")
	}

	if entry.CorrelationID != "test-123" {
		t.Errorf("Expected correlation ID 'test-123', got '%s'", entry.CorrelationID)
	}
}
