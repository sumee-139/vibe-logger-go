package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fladdict/vibe-logger-go/pkg/vibelogger"
)

func main() {
	fmt.Println("=== vibe-logger Configuration Demo ===")
	fmt.Println()

	// Demo 1: Default Configuration
	fmt.Println("--- Demo 1: Default Configuration ---")
	logger1, err := vibelogger.CreateFileLogger("default_demo")
	if err != nil {
		log.Fatalf("Failed to create default logger: %v", err)
	}
	defer logger1.Close()

	logger1.Info("demo", "This uses default configuration")
	fmt.Println()

	// Demo 2: Custom Configuration
	fmt.Println("--- Demo 2: Custom Configuration ---")
	config := &vibelogger.LoggerConfig{
		MaxFileSize:     1024,  // 1KB limit
		AutoSave:        true,
		EnableMemoryLog: true,
		MemoryLogLimit:  3,
		Environment:     "demo",
	}

	logger2, err := vibelogger.CreateFileLoggerWithConfig("custom_demo", config)
	if err != nil {
		log.Fatalf("Failed to create custom logger: %v", err)
	}
	defer logger2.Close()

	logger2.Info("demo", "This uses custom configuration with memory log")
	logger2.Warn("demo", "This is a warning with memory storage")
	logger2.Error("demo", "This is an error with custom settings")

	// Show memory logs
	memoryLogs := logger2.GetMemoryLogs()
	fmt.Printf("Memory logs count: %d\n", len(memoryLogs))
	fmt.Println()

	// Demo 3: Environment Variables
	fmt.Println("--- Demo 3: Environment Variables ---")
	os.Setenv("VIBE_LOG_MAX_FILE_SIZE", "2048")
	os.Setenv("VIBE_LOG_AUTO_SAVE", "false")
	os.Setenv("VIBE_LOG_ENABLE_MEMORY", "true")
	os.Setenv("VIBE_LOG_FILE_PATH", "logs/env_demo.log")
	os.Setenv("VIBE_LOG_ENVIRONMENT", "env_demo")

	envConfig, err := vibelogger.NewConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to create config from environment: %v", err)
	}
	
	logger3, err := vibelogger.CreateFileLoggerWithConfig("env_demo", envConfig)
	if err != nil {
		log.Fatalf("Failed to create env logger: %v", err)
	}
	defer logger3.Close()

	logger3.Info("demo", "This logger was configured from environment variables")
	fmt.Printf("Environment config - MaxFileSize: %d, AutoSave: %t\n", 
		envConfig.MaxFileSize, envConfig.AutoSave)
	fmt.Println()

	// Demo 4: Memory-only Logging
	fmt.Println("--- Demo 4: Memory-only Logging ---")
	memoryConfig := &vibelogger.LoggerConfig{
		AutoSave:        false, // No file writing
		EnableMemoryLog: true,
		MemoryLogLimit:  5,
	}

	logger4 := vibelogger.NewLoggerWithConfig("memory_only", memoryConfig)
	
	for i := 1; i <= 6; i++ {
		logger4.Info("memory_test", fmt.Sprintf("Memory log entry %d", i))
	}

	memoryLogs = logger4.GetMemoryLogs()
	fmt.Printf("Memory-only logs (should be limited to 5): %d entries\n", len(memoryLogs))
	for i, entry := range memoryLogs {
		fmt.Printf("  %d: %s - %s\n", i+1, entry.Operation, entry.Message)
	}

	fmt.Println()
	fmt.Println("=== Demo Complete ===")
	fmt.Println("Check the 'logs' directory for file outputs from demos 1 and 2")

	// Cleanup environment variables
	os.Unsetenv("VIBE_LOG_MAX_FILE_SIZE")
	os.Unsetenv("VIBE_LOG_AUTO_SAVE")
	os.Unsetenv("VIBE_LOG_ENABLE_MEMORY")
	os.Unsetenv("VIBE_LOG_FILE_PATH")
	os.Unsetenv("VIBE_LOG_ENVIRONMENT")
}