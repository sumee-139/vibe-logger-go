package vibelogger

import (
	"errors"
	"testing"
	"time"
)

func TestWithFields(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Test WithFields with multiple field types
	fields := map[string]interface{}{
		"string_field": "test_value",
		"int_field":    123,
		"bool_field":   true,
		"float_field":  45.67,
	}

	err := logger.Info("test_operation", "Test message with fields", WithFields(fields))
	if err != nil {
		t.Fatalf("Failed to log with fields: %v", err)
	}

	// Verify the fields are in the memory log
	logs := logger.GetMemoryLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.Context == nil {
		t.Fatal("Expected context to be set")
	}

	// Check each field
	for key, expectedValue := range fields {
		if actualValue, exists := entry.Context[key]; !exists {
			t.Errorf("Expected field '%s' to exist in context", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected field '%s' to be %v, got %v", key, expectedValue, actualValue)
		}
	}
}

func TestWithError(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Test WithError with a custom error
	testErr := errors.New("test error message")
	err := logger.Error("error_operation", "An error occurred", WithError(testErr))
	if err != nil {
		t.Fatalf("Failed to log with error: %v", err)
	}

	// Verify the error information is in the memory log
	logs := logger.GetMemoryLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.Context == nil {
		t.Fatal("Expected context to be set")
	}

	// Check error field
	if errorStr, exists := entry.Context["error"]; !exists {
		t.Error("Expected 'error' field to exist in context")
	} else if errorStr != "test error message" {
		t.Errorf("Expected error message to be 'test error message', got '%v'", errorStr)
	}

	// Check error_type field
	if errorType, exists := entry.Context["error_type"]; !exists {
		t.Error("Expected 'error_type' field to exist in context")
	} else if errorType != "*errors.errorString" {
		t.Errorf("Expected error type to be '*errors.errorString', got '%v'", errorType)
	}
}

func TestWithUserID(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Test WithUserID
	userID := "user_123456"
	err := logger.Info("user_action", "User performed an action", WithUserID(userID))
	if err != nil {
		t.Fatalf("Failed to log with user ID: %v", err)
	}

	// Verify the user ID is in the memory log
	logs := logger.GetMemoryLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.Context == nil {
		t.Fatal("Expected context to be set")
	}

	// Check user_id field
	if actualUserID, exists := entry.Context["user_id"]; !exists {
		t.Error("Expected 'user_id' field to exist in context")
	} else if actualUserID != userID {
		t.Errorf("Expected user_id to be '%s', got '%v'", userID, actualUserID)
	}
}

func TestWithRequestID(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Test WithRequestID
	requestID := "req_abc123def456"
	err := logger.Info("request_processing", "Processing request", WithRequestID(requestID))
	if err != nil {
		t.Fatalf("Failed to log with request ID: %v", err)
	}

	// Verify the request ID is in the memory log
	logs := logger.GetMemoryLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.Context == nil {
		t.Fatal("Expected context to be set")
	}

	// Check request_id field
	if actualRequestID, exists := entry.Context["request_id"]; !exists {
		t.Error("Expected 'request_id' field to exist in context")
	} else if actualRequestID != requestID {
		t.Errorf("Expected request_id to be '%s', got '%v'", requestID, actualRequestID)
	}
}

func TestMultipleOptions(t *testing.T) {
	// Create a test logger with memory log enabled
	config := &LoggerConfig{
		AutoSave:        false,
		EnableMemoryLog: true,
		MemoryLogLimit:  10,
	}
	logger := NewLoggerWithConfig("test", config)

	// Test multiple options together
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	testErr := errors.New("combined test error")
	userID := "user_789"
	requestID := "req_xyz789"
	correlationID := "corr_123abc"
	duration := 250 * time.Millisecond

	err := logger.Error("complex_operation", "Complex operation with multiple context",
		WithFields(fields),
		WithError(testErr),
		WithUserID(userID),
		WithRequestID(requestID),
		WithCorrelationID(correlationID),
		WithDuration(duration))

	if err != nil {
		t.Fatalf("Failed to log with multiple options: %v", err)
	}

	// Verify all options are applied
	logs := logger.GetMemoryLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	entry := logs[0]
	if entry.Context == nil {
		t.Fatal("Expected context to be set")
	}

	// Check all fields
	expectedContextFields := map[string]interface{}{
		"key1":           "value1",
		"key2":           42,
		"error":          "combined test error",
		"error_type":     "*errors.errorString",
		"user_id":        userID,
		"request_id":     requestID,
		"duration_ms":    int64(250),
		"duration_human": "250ms",
	}

	for key, expectedValue := range expectedContextFields {
		if actualValue, exists := entry.Context[key]; !exists {
			t.Errorf("Expected field '%s' to exist in context", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected field '%s' to be %v, got %v", key, expectedValue, actualValue)
		}
	}

	// Check correlation ID
	if entry.CorrelationID != correlationID {
		t.Errorf("Expected correlation ID to be '%s', got '%s'", correlationID, entry.CorrelationID)
	}
}

func TestWithFieldsNilContext(t *testing.T) {
	// Test WithFields when context is initially nil
	entry := &LogEntry{}
	
	fields := map[string]interface{}{
		"test_key": "test_value",
	}
	
	option := WithFields(fields)
	option(entry)
	
	if entry.Context == nil {
		t.Fatal("Expected context to be initialized")
	}
	
	if value, exists := entry.Context["test_key"]; !exists {
		t.Error("Expected 'test_key' to exist in context")
	} else if value != "test_value" {
		t.Errorf("Expected value to be 'test_value', got '%v'", value)
	}
}

func TestWithErrorNilContext(t *testing.T) {
	// Test WithError when context is initially nil
	entry := &LogEntry{}
	
	testErr := errors.New("nil context test error")
	option := WithError(testErr)
	option(entry)
	
	if entry.Context == nil {
		t.Fatal("Expected context to be initialized")
	}
	
	if errorStr, exists := entry.Context["error"]; !exists {
		t.Error("Expected 'error' field to exist in context")
	} else if errorStr != "nil context test error" {
		t.Errorf("Expected error message to be 'nil context test error', got '%v'", errorStr)
	}
}

func TestWithUserIDNilContext(t *testing.T) {
	// Test WithUserID when context is initially nil
	entry := &LogEntry{}
	
	userID := "test_user_nil_context"
	option := WithUserID(userID)
	option(entry)
	
	if entry.Context == nil {
		t.Fatal("Expected context to be initialized")
	}
	
	if actualUserID, exists := entry.Context["user_id"]; !exists {
		t.Error("Expected 'user_id' field to exist in context")
	} else if actualUserID != userID {
		t.Errorf("Expected user_id to be '%s', got '%v'", userID, actualUserID)
	}
}

func TestWithRequestIDNilContext(t *testing.T) {
	// Test WithRequestID when context is initially nil
	entry := &LogEntry{}
	
	requestID := "test_request_nil_context"
	option := WithRequestID(requestID)
	option(entry)
	
	if entry.Context == nil {
		t.Fatal("Expected context to be initialized")
	}
	
	if actualRequestID, exists := entry.Context["request_id"]; !exists {
		t.Error("Expected 'request_id' field to exist in context")
	} else if actualRequestID != requestID {
		t.Errorf("Expected request_id to be '%s', got '%v'", requestID, actualRequestID)
	}
}