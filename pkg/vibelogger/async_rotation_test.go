package vibelogger

import (
	"os"
	"testing"
	"time"
)

func TestAsyncRotation(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     200, // Small for testing
		RotationEnabled: true,
		MaxRotatedFiles: 3,
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/async_rotation_test.log"

	logger, err := CreateFileLoggerWithConfig("async_rotation_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write logs to trigger rotation
	for i := 0; i < 10; i++ {
		err := logger.Info("async_test", "This is a test message for async rotation testing")
		if err != nil {
			t.Fatalf("Failed to write log entry %d: %v", i, err)
		}
	}

	// Test async rotation
	errCh := logger.ForceRotationAsync()
	
	// Should return immediately without blocking
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Async rotation failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Async rotation took too long")
	}

	// Verify rotation occurred
	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) == 0 {
		t.Error("Expected at least one rotated file after async rotation")
	}

	// Test writing after async rotation
	err = logger.Info("async_test", "Log after async rotation")
	if err != nil {
		t.Fatalf("Failed to write after async rotation: %v", err)
	}
}

func TestAsyncRotationDisabled(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     200,
		RotationEnabled: true,
		MaxRotatedFiles: 3,
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/async_disabled_test.log"

	logger, err := CreateFileLoggerWithConfig("async_disabled_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Disable async rotation
	logger.SetAsyncRotation(false)

	// Test async rotation with disabled async
	errCh := logger.ForceRotationAsync()
	
	// Should still work (fallback to sync)
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Async rotation failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Async rotation took too long")
	}
}

func TestPerformanceImprovement(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     1000, // Small for testing
		RotationEnabled: true,
		MaxRotatedFiles: 5,
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/performance_test.log"

	logger, err := CreateFileLoggerWithConfig("performance_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Test that file size caching is working
	start := time.Now()
	
	// Write many small logs to test cached size performance
	for i := 0; i < 100; i++ {
		err := logger.Info("perf_test", "Short log message for performance testing")
		if err != nil {
			t.Fatalf("Failed to write log entry %d: %v", i, err)
		}
	}
	
	duration := time.Since(start)
	
	// Should complete quickly with cached sizing
	if duration > 5*time.Second {
		t.Errorf("Performance test took too long: %v", duration)
	}
	
	t.Logf("Performance test completed in: %v", duration)
}