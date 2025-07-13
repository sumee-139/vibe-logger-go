package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sumee-139/vibe-logger-go/pkg/vibelogger"
)

func main() {
	fmt.Println("=== vibe-logger Log Rotation Demo ===")
	fmt.Println()

	// Cleanup any existing test files
	os.RemoveAll("demo_logs")

	// Demo 1: Basic Log Rotation
	fmt.Println("--- Demo 1: Basic Log Rotation ---")
	config := &vibelogger.LoggerConfig{
		MaxFileSize:     500, // Very small for demo purposes (500 bytes)
		RotationEnabled: true,
		MaxRotatedFiles: 3, // Keep 3 rotated files
		AutoSave:        true,
		FilePath:        "demo_logs/rotation_demo.log",
	}

	logger, err := vibelogger.CreateFileLoggerWithConfig("rotation_demo", config)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Generate enough logs to trigger multiple rotations
	fmt.Println("Writing logs to trigger rotation...")
	for i := 1; i <= 20; i++ {
		err := logger.Info("demo", fmt.Sprintf("This is log entry #%d which should eventually trigger log rotation when the file gets large enough", i))
		if err != nil {
			log.Printf("Failed to write log entry %d: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond) // Small delay for different timestamps
	}

	// Show rotation results
	rotatedFiles := logger.GetRotatedFiles()
	fmt.Printf("Number of rotated files: %d\n", len(rotatedFiles))
	for i, file := range rotatedFiles {
		fmt.Printf("  Rotated file %d: %s\n", i+1, file)
	}
	fmt.Println()

	// Demo 2: Manual Rotation
	fmt.Println("--- Demo 2: Manual Rotation ---")
	fmt.Println("Current file size before manual rotation...")

	// Add a few more logs
	logger.Info("manual", "Log before manual rotation")
	logger.Info("manual", "Another log before manual rotation")

	// Force manual rotation
	fmt.Println("Forcing manual rotation...")
	err = logger.ForceRotation()
	if err != nil {
		log.Printf("Failed to force rotation: %v", err)
	} else {
		fmt.Println("Manual rotation completed successfully!")
	}

	// Add logs after rotation
	logger.Info("manual", "Log after manual rotation")
	logger.Info("manual", "Final log after manual rotation")

	// Show final rotation results
	finalRotatedFiles := logger.GetRotatedFiles()
	fmt.Printf("Final number of rotated files: %d\n", len(finalRotatedFiles))
	fmt.Println()

	// Demo 3: Environment Variable Configuration
	fmt.Println("--- Demo 3: Environment Variable Configuration ---")

	// Set environment variables for rotation
	os.Setenv("VIBE_LOG_MAX_FILE_SIZE", "300")
	os.Setenv("VIBE_LOG_ROTATION_ENABLED", "true")
	os.Setenv("VIBE_LOG_MAX_ROTATED_FILES", "2")
	os.Setenv("VIBE_LOG_FILE_PATH", "demo_logs/env_rotation.log")

	envConfig, err := vibelogger.NewConfigFromEnvironment()
	if err != nil {
		log.Fatalf("Failed to create config from environment: %v", err)
	}

	envLogger, err := vibelogger.CreateFileLoggerWithConfig("env_rotation", envConfig)
	if err != nil {
		log.Fatalf("Failed to create env logger: %v", err)
	}
	defer envLogger.Close()

	fmt.Printf("Environment config - MaxFileSize: %d, RotationEnabled: %t, MaxRotatedFiles: %d\n",
		envConfig.MaxFileSize, envConfig.RotationEnabled, envConfig.MaxRotatedFiles)

	// Generate logs to trigger rotation
	for i := 1; i <= 15; i++ {
		envLogger.Info("env_demo", fmt.Sprintf("Environment configured log #%d", i))
		time.Sleep(5 * time.Millisecond)
	}

	envRotatedFiles := envLogger.GetRotatedFiles()
	fmt.Printf("Environment demo rotated files: %d\n", len(envRotatedFiles))
	fmt.Println()

	// Demo 4: Rotation Disabled
	fmt.Println("--- Demo 4: Rotation Disabled ---")

	disabledConfig := &vibelogger.LoggerConfig{
		MaxFileSize:     100, // Small size but rotation disabled
		RotationEnabled: false,
		AutoSave:        true,
		FilePath:        "demo_logs/no_rotation.log",
	}

	noRotationLogger, err := vibelogger.CreateFileLoggerWithConfig("no_rotation", disabledConfig)
	if err != nil {
		log.Fatalf("Failed to create no-rotation logger: %v", err)
	}
	defer noRotationLogger.Close()

	fmt.Println("Writing logs with rotation disabled...")
	for i := 1; i <= 10; i++ {
		noRotationLogger.Info("no_rotation", fmt.Sprintf("Log #%d without rotation", i))
	}

	noRotationFiles := noRotationLogger.GetRotatedFiles()
	fmt.Printf("No-rotation demo rotated files: %d (should be 0)\n", len(noRotationFiles))
	fmt.Println()

	// Show directory contents
	fmt.Println("--- Final Directory Contents ---")
	files, err := os.ReadDir("demo_logs")
	if err != nil {
		log.Printf("Failed to read demo directory: %v", err)
	} else {
		fmt.Printf("Files in demo_logs directory:\n")
		for _, file := range files {
			if !file.IsDir() {
				info, _ := file.Info()
				fmt.Printf("  %s (%d bytes)\n", file.Name(), info.Size())
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Demo Complete ===")
	fmt.Println("Check the 'demo_logs' directory to see the rotation results!")
	fmt.Println("Notice how files are automatically rotated with timestamps and old files are cleaned up.")

	// Cleanup environment variables
	os.Unsetenv("VIBE_LOG_MAX_FILE_SIZE")
	os.Unsetenv("VIBE_LOG_ROTATION_ENABLED")
	os.Unsetenv("VIBE_LOG_MAX_ROTATED_FILES")
	os.Unsetenv("VIBE_LOG_FILE_PATH")
}
