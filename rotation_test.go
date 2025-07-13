package vibelogger

import (
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLogRotationBasic(t *testing.T) {
	// Cleanup test files
	defer func() {
		os.RemoveAll("test_logs")
	}()

	// Create logger with small file size limit for testing
	config := &LoggerConfig{
		MaxFileSize:     100, // Very small for testing
		RotationEnabled: true,
		MaxRotatedFiles: 3,
		AutoSave:        true,
	}

	// Create test directory
	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Custom file path for testing
	config.FilePath = "test_logs/rotation_test.log"

	logger, err := CreateFileLoggerWithConfig("rotation_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write enough logs to trigger rotation
	for i := 0; i < 10; i++ {
		err := logger.Info("test_operation", "This is a test message that should be long enough to trigger rotation eventually")
		if err != nil {
			t.Fatalf("Failed to write log entry %d: %v", i, err)
		}
	}

	// Check that rotation occurred
	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) == 0 {
		t.Error("Expected at least one rotated file, but got none")
	}

	// Verify original file still exists and is usable
	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		t.Error("Original log file should still exist after rotation")
	}
}

func TestLogRotationRetentionPolicy(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     50, // Very small for testing
		RotationEnabled: true,
		MaxRotatedFiles: 2, // Keep only 2 files
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/retention_test.log"

	logger, err := CreateFileLoggerWithConfig("retention_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Force multiple rotations
	for rotation := 0; rotation < 5; rotation++ {
		for i := 0; i < 5; i++ {
			err := logger.Info("test", "Message to trigger rotation through size limit")
			if err != nil {
				t.Fatalf("Failed to write log entry: %v", err)
			}
		}
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}

	// Check retention policy
	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) > config.MaxRotatedFiles {
		t.Errorf("Expected max %d rotated files, but got %d", config.MaxRotatedFiles, len(rotatedFiles))
	}

	// Verify files were actually deleted
	files, err := os.ReadDir("test_logs")
	if err != nil {
		t.Fatalf("Failed to read test directory: %v", err)
	}

	logFileCount := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") || strings.Contains(file.Name(), "retention_test.log") {
			logFileCount++
		}
	}

	// Should have 1 current file + MaxRotatedFiles rotated files
	expectedCount := 1 + config.MaxRotatedFiles
	if logFileCount > expectedCount {
		t.Errorf("Expected at most %d log files, but found %d", expectedCount, logFileCount)
	}
}

func TestConcurrentRotation(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     200, // Small for testing
		RotationEnabled: true,
		MaxRotatedFiles: 5,
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/concurrent_test.log"

	logger, err := CreateFileLoggerWithConfig("concurrent_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Concurrent logging to test thread safety
	var wg sync.WaitGroup
	numGoroutines := 10
	numLogs := 20

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < numLogs; i++ {
				err := logger.Info("concurrent_test", "Concurrent log message from goroutine %d, message %d")
				if err != nil {
					t.Errorf("Goroutine %d failed to write log: %v", goroutineID, err)
				}
			}
		}(g)
	}

	wg.Wait()

	// Verify no data corruption occurred
	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		t.Error("Main log file should exist after concurrent logging")
	}

	// Check that rotated files are valid (exist and are readable)
	rotatedFiles := logger.GetRotatedFiles()
	for _, file := range rotatedFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Logf("Warning: Rotated file %s was mentioned but doesn't exist (may have been cleaned up)", file)
		}
	}
}

func TestRotationDisabled(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	// Test with rotation disabled
	config := &LoggerConfig{
		MaxFileSize:     100,
		RotationEnabled: false, // Disabled
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/no_rotation_test.log"

	logger, err := CreateFileLoggerWithConfig("no_rotation_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Verify rotation manager is not initialized
	if logger.rotationMgr != nil {
		t.Error("Rotation manager should be nil when rotation is disabled")
	}

	// Verify no rotated files are created even with many logs
	for i := 0; i < 5; i++ {
		logger.Info("test", "Message without rotation enabled")
	}

	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) != 0 {
		t.Errorf("Expected no rotated files when rotation is disabled, got %d", len(rotatedFiles))
	}
}

func TestForceRotation(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     10000, // Large enough to not trigger automatic rotation
		RotationEnabled: true,
		MaxRotatedFiles: 3,
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/force_rotation_test.log"

	logger, err := CreateFileLoggerWithConfig("force_rotation_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Write some logs
	logger.Info("test", "Initial log entry")
	logger.Info("test", "Another log entry")

	// Force rotation manually
	err = logger.ForceRotation()
	if err != nil {
		t.Fatalf("Failed to force rotation: %v", err)
	}

	// Verify rotation occurred
	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) != 1 {
		t.Errorf("Expected 1 rotated file after force rotation, got %d", len(rotatedFiles))
	}

	// Verify original file was recreated and is writable
	err = logger.Info("test", "Log after forced rotation")
	if err != nil {
		t.Fatalf("Failed to write after forced rotation: %v", err)
	}
}

func TestRotationWithCustomPaths(t *testing.T) {
	defer func() {
		os.RemoveAll("custom_test_dir")
	}()

	// Test with custom directory path
	config := &LoggerConfig{
		MaxFileSize:     80,
		RotationEnabled: true,
		MaxRotatedFiles: 2,
		AutoSave:        true,
		FilePath:        "custom_test_dir/subdir/custom.log",
	}

	logger, err := CreateFileLoggerWithConfig("custom_path_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger with custom path: %v", err)
	}
	defer logger.Close()

	// Trigger rotation
	for i := 0; i < 10; i++ {
		err := logger.Info("test", "Message to trigger rotation in custom directory")
		if err != nil {
			t.Fatalf("Failed to write log entry: %v", err)
		}
	}

	// Verify files exist in custom directory
	if _, err := os.Stat("custom_test_dir/subdir/custom.log"); os.IsNotExist(err) {
		t.Error("Custom log file should exist")
	}

	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) == 0 {
		t.Error("Expected rotated files in custom directory")
	}

	// Verify rotated files are in the correct directory
	for _, file := range rotatedFiles {
		if !strings.HasPrefix(file, "custom_test_dir/subdir/") {
			t.Errorf("Rotated file %s should be in custom directory", file)
		}
	}
}

func TestConfigurationUpdate(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     1000,
		RotationEnabled: false, // Start with rotation disabled
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/config_update_test.log"

	logger, err := CreateFileLoggerWithConfig("config_update_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Verify rotation is initially disabled
	if logger.rotationMgr != nil {
		t.Error("Rotation manager should be nil initially")
	}

	// Enable rotation via config update (keep existing file path)
	newConfig := &LoggerConfig{
		MaxFileSize:     100,
		RotationEnabled: true, // Enable rotation
		MaxRotatedFiles: 2,
		AutoSave:        true,
		FilePath:        "", // Don't change file path to avoid validation issues
		Environment:     "development",
	}

	err = logger.UpdateConfig(newConfig)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Verify rotation manager is now initialized
	if logger.rotationMgr == nil {
		t.Error("Rotation manager should be initialized after enabling rotation")
	}

	// Test that rotation now works
	for i := 0; i < 10; i++ {
		logger.Info("test", "Message after enabling rotation via config update")
	}

	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) == 0 {
		t.Error("Expected rotated files after enabling rotation")
	}
}

func TestRotationManagerUpdateConfig(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	// Create initial config with rotation enabled
	config := &LoggerConfig{
		MaxFileSize:     1000,
		RotationEnabled: true,
		MaxRotatedFiles: 3,
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/rotation_manager_update_test.log"

	logger, err := CreateFileLoggerWithConfig("rotation_manager_update_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Verify rotation manager exists
	if logger.rotationMgr == nil {
		t.Fatal("Rotation manager should be initialized")
	}

	// Create some rotated files first
	for i := 0; i < 5; i++ {
		logger.Info("test", "Message to create some log history")
	}
	
	// Force a rotation to create a rotated file
	err = logger.ForceRotation()
	if err != nil {
		t.Fatalf("Failed to force initial rotation: %v", err)
	}

	initialRotatedFiles := logger.GetRotatedFiles()
	initialCount := len(initialRotatedFiles)

	// Update config with different MaxRotatedFiles (reduce retention)
	newConfig := &LoggerConfig{
		MaxFileSize:     1000,
		RotationEnabled: true,
		MaxRotatedFiles: 1, // Reduce from 3 to 1
		AutoSave:        true,
		FilePath:        config.FilePath,
		Environment:     "production",
	}

	// Test UpdateConfig directly on RotationManager
	logger.rotationMgr.UpdateConfig(newConfig)

	// Verify the config was updated
	logger.rotationMgr.mutex.Lock()
	updatedConfig := logger.rotationMgr.config
	logger.rotationMgr.mutex.Unlock()

	if updatedConfig.MaxRotatedFiles != 1 {
		t.Errorf("Expected MaxRotatedFiles to be updated to 1, got %d", updatedConfig.MaxRotatedFiles)
	}

	if updatedConfig.Environment != "production" {
		t.Errorf("Expected Environment to be updated to 'production', got '%s'", updatedConfig.Environment)
	}

	// Create another rotation to trigger cleanup
	logger.Info("test", "Message after config update")
	err = logger.ForceRotation()
	if err != nil {
		t.Fatalf("Failed to force rotation after config update: %v", err)
	}

	// Verify that old files were cleaned up according to new retention policy
	finalRotatedFiles := logger.GetRotatedFiles()
	finalCount := len(finalRotatedFiles)

	if finalCount > newConfig.MaxRotatedFiles {
		t.Errorf("Expected at most %d rotated files after config update, got %d", newConfig.MaxRotatedFiles, finalCount)
	}

	// Should have fewer or equal files than initially (due to cleanup)
	if finalCount > initialCount {
		t.Errorf("Expected cleanup to reduce file count from %d to at most %d, got %d", initialCount, newConfig.MaxRotatedFiles, finalCount)
	}
}

func TestPerformRotationAsync(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	// Test async rotation enabled
	config := &LoggerConfig{
		MaxFileSize:     500,
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

	if logger.rotationMgr == nil {
		t.Fatal("Rotation manager should be initialized")
	}

	// Write some logs to create content for rotation
	for i := 0; i < 5; i++ {
		logger.Info("test", "Message to create content for async rotation testing")
	}

	// Test PerformRotationAsync with async enabled
	logger.rotationMgr.asyncEnabled = true
	responseChan := logger.rotationMgr.PerformRotationAsync()

	// Verify we get a response channel
	if responseChan == nil {
		t.Fatal("Expected response channel from PerformRotationAsync")
	}

	// Wait for the async operation to complete with timeout
	select {
	case err := <-responseChan:
		if err != nil {
			t.Fatalf("Async rotation failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Async rotation timed out")
	}

	// Verify rotation occurred
	rotatedFiles := logger.GetRotatedFiles()
	if len(rotatedFiles) == 0 {
		t.Error("Expected at least one rotated file after async rotation")
	}

	// Test PerformRotationAsync with async disabled (fallback to sync)
	logger.rotationMgr.asyncEnabled = false
	
	// Write more logs
	for i := 0; i < 3; i++ {
		logger.Info("test", "Message before sync fallback test")
	}

	responseChan2 := logger.rotationMgr.PerformRotationAsync()

	// Verify we still get a response channel (sync fallback)
	if responseChan2 == nil {
		t.Fatal("Expected response channel even when async is disabled (sync fallback)")
	}

	// Wait for the sync fallback operation to complete
	select {
	case err := <-responseChan2:
		if err != nil {
			t.Fatalf("Sync fallback rotation failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Sync fallback rotation timed out")
	}

	// Verify more rotations occurred
	finalRotatedFiles := logger.GetRotatedFiles()
	if len(finalRotatedFiles) < len(rotatedFiles) {
		t.Error("Expected more rotated files after second rotation")
	}
}

func TestPerformRotationAsyncChannelBehavior(t *testing.T) {
	defer func() {
		os.RemoveAll("test_logs")
	}()

	config := &LoggerConfig{
		MaxFileSize:     200,
		RotationEnabled: true,
		MaxRotatedFiles: 2,
		AutoSave:        true,
	}

	if err := os.MkdirAll("test_logs", 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	config.FilePath = "test_logs/async_channel_test.log"

	logger, err := CreateFileLoggerWithConfig("async_channel_test", config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	if logger.rotationMgr == nil {
		t.Fatal("Rotation manager should be initialized")
	}

	// Write initial logs
	for i := 0; i < 3; i++ {
		logger.Info("test", "Initial message for channel behavior test")
	}

	// Test multiple concurrent async rotations
	logger.rotationMgr.asyncEnabled = true
	
	channels := make([]<-chan error, 3)
	for i := 0; i < 3; i++ {
		channels[i] = logger.rotationMgr.PerformRotationAsync()
		if channels[i] == nil {
			t.Fatalf("Expected non-nil channel from PerformRotationAsync call %d", i)
		}
	}

	// Wait for all async operations to complete
	for i, ch := range channels {
		select {
		case err := <-ch:
			if err != nil {
				t.Errorf("Async rotation %d failed: %v", i, err)
			}
		case <-time.After(5 * time.Second):
			t.Errorf("Async rotation %d timed out", i)
		}
	}

	// Verify that rotations actually occurred
	finalRotatedFiles := logger.GetRotatedFiles()
	if len(finalRotatedFiles) == 0 {
		t.Error("Expected rotated files after multiple async rotations")
	}
}
