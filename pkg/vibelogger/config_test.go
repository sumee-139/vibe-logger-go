package vibelogger

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.MaxFileSize != 10*1024*1024 {
		t.Errorf("Expected MaxFileSize to be 10MB, got %d", config.MaxFileSize)
	}

	if !config.AutoSave {
		t.Error("Expected AutoSave to be true by default")
	}

	if config.EnableMemoryLog {
		t.Error("Expected EnableMemoryLog to be false by default")
	}

	if config.MemoryLogLimit != 1000 {
		t.Errorf("Expected MemoryLogLimit to be 1000, got %d", config.MemoryLogLimit)
	}

	if config.Environment != "development" {
		t.Errorf("Expected Environment to be 'development', got '%s'", config.Environment)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	// Set test environment variables
	os.Setenv("VIBE_LOG_MAX_FILE_SIZE", "5242880") // 5MB
	os.Setenv("VIBE_LOG_AUTO_SAVE", "false")
	os.Setenv("VIBE_LOG_ENABLE_MEMORY", "true")
	os.Setenv("VIBE_LOG_MEMORY_LIMIT", "500")
	os.Setenv("VIBE_LOG_FILE_PATH", "logs/test.log") // Use safe relative path
	os.Setenv("VIBE_LOG_ENVIRONMENT", "test")

	defer func() {
		// Clean up environment variables
		os.Unsetenv("VIBE_LOG_MAX_FILE_SIZE")
		os.Unsetenv("VIBE_LOG_AUTO_SAVE")
		os.Unsetenv("VIBE_LOG_ENABLE_MEMORY")
		os.Unsetenv("VIBE_LOG_MEMORY_LIMIT")
		os.Unsetenv("VIBE_LOG_FILE_PATH")
		os.Unsetenv("VIBE_LOG_ENVIRONMENT")
	}()

	config, err := NewConfigFromEnvironment()
	if err != nil {
		t.Fatalf("Failed to create config from environment: %v", err)
	}

	if config.MaxFileSize != 5242880 {
		t.Errorf("Expected MaxFileSize to be 5242880, got %d", config.MaxFileSize)
	}

	if config.AutoSave {
		t.Error("Expected AutoSave to be false")
	}

	if !config.EnableMemoryLog {
		t.Error("Expected EnableMemoryLog to be true")
	}

	if config.MemoryLogLimit != 500 {
		t.Errorf("Expected MemoryLogLimit to be 500, got %d", config.MemoryLogLimit)
	}

	if config.FilePath != "logs/test.log" {
		t.Errorf("Expected FilePath to be 'logs/test.log', got '%s'", config.FilePath)
	}

	if config.Environment != "test" {
		t.Errorf("Expected Environment to be 'test', got '%s'", config.Environment)
	}
}

func TestConfigValidate(t *testing.T) {
	config := &LoggerConfig{
		MaxFileSize:    -100,
		MemoryLogLimit: -50,
		Environment:    "",
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate should not return error, got: %v", err)
	}

	if config.MaxFileSize != 0 {
		t.Errorf("Expected negative MaxFileSize to be set to 0, got %d", config.MaxFileSize)
	}

	if config.MemoryLogLimit != 0 {
		t.Errorf("Expected negative MemoryLogLimit to be set to 0, got %d", config.MemoryLogLimit)
	}

	if config.Environment != "development" {
		t.Errorf("Expected empty Environment to be set to 'development', got '%s'", config.Environment)
	}
}

func TestCreateFileLoggerWithConfig(t *testing.T) {
	defer func() {
		os.RemoveAll("logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     1024, // 1KB limit for testing
		AutoSave:        true,
		EnableMemoryLog: true,
		MemoryLogLimit:  5,
	}

	logger, err := CreateFileLoggerWithConfig("test_config", config)
	if err != nil {
		t.Fatalf("Failed to create logger with config: %v", err)
	}
	defer logger.Close()

	if logger.config.MaxFileSize != 1024 {
		t.Errorf("Expected logger MaxFileSize to be 1024, got %d", logger.config.MaxFileSize)
	}

	if !logger.config.EnableMemoryLog {
		t.Error("Expected logger EnableMemoryLog to be true")
	}
}

func TestMemoryLogFunctionality(t *testing.T) {
	config := &LoggerConfig{
		AutoSave:        false, // Disable file writing for this test
		EnableMemoryLog: true,
		MemoryLogLimit:  3,
	}

	logger := NewLoggerWithConfig("test_memory", config)

	// Log some entries
	logger.Info("op1", "message1")
	logger.Info("op2", "message2")
	logger.Info("op3", "message3")

	memoryLogs := logger.GetMemoryLogs()
	if len(memoryLogs) != 3 {
		t.Errorf("Expected 3 memory logs, got %d", len(memoryLogs))
	}

	// Add one more to trigger limit
	logger.Info("op4", "message4")

	memoryLogs = logger.GetMemoryLogs()
	if len(memoryLogs) != 3 {
		t.Errorf("Expected memory logs to be limited to 3, got %d", len(memoryLogs))
	}

	// Check that oldest log was removed
	if memoryLogs[0].Operation != "op2" {
		t.Errorf("Expected first log operation to be 'op2', got '%s'", memoryLogs[0].Operation)
	}

	// Test clear functionality
	logger.ClearMemoryLogs()
	memoryLogs = logger.GetMemoryLogs()
	if len(memoryLogs) != 0 {
		t.Errorf("Expected memory logs to be cleared, got %d entries", len(memoryLogs))
	}
}

func TestAutoSaveToggle(t *testing.T) {
	defer func() {
		os.RemoveAll("logs")
	}()

	// Test with AutoSave disabled
	config := &LoggerConfig{
		AutoSave: false,
	}

	logger, err := CreateFileLoggerWithConfig("test_autosave", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Log a message - should not write to file
	err = logger.Info("test_op", "test message")
	if err != nil {
		t.Errorf("Logging should not error even with AutoSave disabled: %v", err)
	}

	// Check file size should still be 0 or minimal
	if logger.currentSize > 0 {
		t.Errorf("Expected file size to be 0 with AutoSave disabled, got %d", logger.currentSize)
	}
}

// Security Tests

func TestPathTraversalPrevention(t *testing.T) {
	// Test various path traversal attack vectors
	dangerousPaths := []string{
		"../../../etc/passwd",
		"..\\..\\windows\\system32\\config\\sam",
		"/etc/passwd",
		"logs/../../../sensitive.log",
		"logs\\..\\..\\sensitive.log",
		"logs/../../home/user/.ssh/id_rsa",
	}

	for _, path := range dangerousPaths {
		config := &LoggerConfig{
			FilePath: path,
		}

		err := config.Validate()
		if err == nil {
			t.Errorf("Expected validation to fail for dangerous path: %s", path)
		}
	}
}

func TestEnvironmentVariableInjection(t *testing.T) {
	// Test malicious environment variable values
	defer func() {
		os.Unsetenv("VIBE_LOG_MAX_FILE_SIZE")
		os.Unsetenv("VIBE_LOG_MEMORY_LIMIT")
		os.Unsetenv("VIBE_LOG_FILE_PATH")
		os.Unsetenv("VIBE_LOG_ENVIRONMENT")
	}()

	// Test invalid file size values
	os.Setenv("VIBE_LOG_MAX_FILE_SIZE", "999999999999999999999") // Too large
	_, err := NewConfigFromEnvironment()
	if err == nil {
		t.Error("Expected error for oversized VIBE_LOG_MAX_FILE_SIZE")
	}

	// Test negative values
	os.Setenv("VIBE_LOG_MAX_FILE_SIZE", "-1000")
	_, err = NewConfigFromEnvironment()
	if err == nil {
		t.Error("Expected error for negative VIBE_LOG_MAX_FILE_SIZE")
	}

	// Test invalid memory limit
	os.Setenv("VIBE_LOG_MEMORY_LIMIT", "99999999") // Too large
	_, err = NewConfigFromEnvironment()
	if err == nil {
		t.Error("Expected error for oversized VIBE_LOG_MEMORY_LIMIT")
	}

	// Test path traversal in environment
	os.Setenv("VIBE_LOG_FILE_PATH", "../../../etc/passwd")
	_, err = NewConfigFromEnvironment()
	if err == nil {
		t.Error("Expected error for path traversal in VIBE_LOG_FILE_PATH")
	}

	// Test invalid environment names
	invalidEnvNames := []string{
		"env with spaces",
		"env;with;semicolons",
		"env|with|pipes",
		"env$(command)",
		"env`command`",
		"env&command&",
	}

	for _, envName := range invalidEnvNames {
		os.Setenv("VIBE_LOG_ENVIRONMENT", envName)
		_, err = NewConfigFromEnvironment()
		if err == nil {
			t.Errorf("Expected error for invalid environment name: %s", envName)
		}
	}
}

func TestResourceLimits(t *testing.T) {
	// Test that resource limits are enforced
	config := &LoggerConfig{
		MaxFileSize:    MaxFileSizeLimit + 1,  // Exceed limit
		MemoryLogLimit: MaxMemoryLogLimit + 1, // Exceed limit
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected validation to fail for configs exceeding resource limits")
	}
}

func TestValidEnvironmentNames(t *testing.T) {
	// Test that valid environment names are accepted
	validEnvNames := []string{
		"development",
		"test",
		"production",
		"staging",
		"dev-1",
		"prod_v2",
		"test.env",
		"DEV",
		"PROD",
	}

	for _, envName := range validEnvNames {
		if !isValidEnvironmentName(envName) {
			t.Errorf("Expected environment name to be valid: %s", envName)
		}
	}
}

func TestSecurePathValidation(t *testing.T) {
	// Test that safe paths are accepted
	safePaths := []string{
		"logs/app.log",
		"./logs/app.log",
		"logs/subdir/app.log",
		"/tmp/app.log",
		"/var/log/app.log",
		"", // Empty path should be valid (uses default)
	}

	for _, path := range safePaths {
		config := &LoggerConfig{
			FilePath: path,
		}

		err := config.validateFilePath()
		if err != nil {
			t.Errorf("Expected path to be valid: %s, got error: %v", path, err)
		}
	}
}
