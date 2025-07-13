package vibelogger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// RotationManager handles log file rotation and cleanup
type RotationManager struct {
	logger       *Logger
	config       *LoggerConfig
	basePath     string
	mutex        sync.Mutex
	rotatedFiles []string
}

// NewRotationManager creates a new rotation manager for the given logger
func NewRotationManager(logger *Logger, config *LoggerConfig, basePath string) *RotationManager {
	rm := &RotationManager{
		logger:   logger,
		config:   config,
		basePath: basePath,
	}

	// Initialize list of existing rotated files
	rm.scanExistingRotatedFiles()

	return rm
}

// ShouldRotate checks if rotation is needed for the given entry size
func (rm *RotationManager) ShouldRotate(newEntrySize int64) bool {
	if !rm.config.RotationEnabled {
		return false
	}

	if rm.config.MaxFileSize <= 0 {
		return false // Unlimited file size
	}

	return rm.logger.currentSize+newEntrySize > rm.config.MaxFileSize
}

// PerformRotation rotates the current log file and creates a new one
func (rm *RotationManager) PerformRotation() error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Close current file
	if rm.logger.file != nil {
		if err := rm.logger.file.Close(); err != nil {
			return fmt.Errorf("failed to close current log file: %w", err)
		}
	}

	// Generate rotated file name with timestamp
	timestamp := time.Now().Format("20060102_150405")
	rotatedPath := fmt.Sprintf("%s.%s", rm.basePath, timestamp)

	// Rename current file to rotated name
	if err := os.Rename(rm.basePath, rotatedPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Add to rotated files list
	rm.rotatedFiles = append(rm.rotatedFiles, rotatedPath)

	// Clean up old files if needed
	if err := rm.cleanupOldFiles(); err != nil {
		// Log warning but don't fail rotation
		fmt.Printf("Warning: failed to cleanup old files: %v\n", err)
	}

	// Create new log file
	newFile, err := os.OpenFile(rm.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	// Update logger with new file
	rm.logger.file = newFile
	rm.logger.currentSize = 0

	return nil
}

// scanExistingRotatedFiles scans for existing rotated files matching the pattern
func (rm *RotationManager) scanExistingRotatedFiles() {
	baseDir := filepath.Dir(rm.basePath)
	baseName := filepath.Base(rm.basePath)

	files, err := os.ReadDir(baseDir)
	if err != nil {
		return // Directory doesn't exist or can't read
	}

	var rotatedFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		// Check if it's a rotated file (baseName.timestamp)
		if strings.HasPrefix(name, baseName+".") {
			rotatedFiles = append(rotatedFiles, filepath.Join(baseDir, name))
		}
	}

	// Sort by modification time (newest first)
	sort.Slice(rotatedFiles, func(i, j int) bool {
		infoI, errI := os.Stat(rotatedFiles[i])
		infoJ, errJ := os.Stat(rotatedFiles[j])
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().After(infoJ.ModTime())
	})

	rm.rotatedFiles = rotatedFiles
}

// cleanupOldFiles removes old rotated files based on retention policy
func (rm *RotationManager) cleanupOldFiles() error {
	if rm.config.MaxRotatedFiles <= 0 {
		return nil // Keep all files
	}

	// Sort files by modification time (newest first)
	sort.Slice(rm.rotatedFiles, func(i, j int) bool {
		infoI, errI := os.Stat(rm.rotatedFiles[i])
		infoJ, errJ := os.Stat(rm.rotatedFiles[j])
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Remove files beyond the retention limit
	if len(rm.rotatedFiles) > rm.config.MaxRotatedFiles {
		filesToDelete := rm.rotatedFiles[rm.config.MaxRotatedFiles:]

		for _, file := range filesToDelete {
			if err := os.Remove(file); err != nil {
				return fmt.Errorf("failed to remove old rotated file %s: %w", file, err)
			}
		}

		// Update the list
		rm.rotatedFiles = rm.rotatedFiles[:rm.config.MaxRotatedFiles]
	}

	return nil
}

// GetRotatedFiles returns the list of current rotated files
func (rm *RotationManager) GetRotatedFiles() []string {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Return a copy to prevent external modification
	files := make([]string, len(rm.rotatedFiles))
	copy(files, rm.rotatedFiles)
	return files
}

// UpdateConfig updates the rotation manager configuration
func (rm *RotationManager) UpdateConfig(config *LoggerConfig) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	rm.config = config

	// Clean up files if retention policy changed
	if err := rm.cleanupOldFiles(); err != nil {
		fmt.Printf("Warning: failed to cleanup files after config update: %v\n", err)
	}
}
