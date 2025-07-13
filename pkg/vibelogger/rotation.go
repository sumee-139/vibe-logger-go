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

// rotationRequest は非同期ローテーション要求を表す
type rotationRequest struct {
	force    bool           // 強制ローテーションかどうか
	response chan error     // 結果を返すチャネル
}

// RotationManager handles log file rotation and cleanup
type RotationManager struct {
	logger              *Logger
	config              *LoggerConfig
	basePath            string
	mutex               sync.Mutex
	rotatedFiles        []string
	cachedFileSize      int64     // Cached file size for performance
	lastSizeSync        time.Time // Last time we synced with actual file size
	sizeSyncInterval    time.Duration // How often to sync cached size with disk
	pendingRotation     bool      // Flag to prevent duplicate rotations
	asyncRotationChan   chan rotationRequest // Channel for async rotation requests
	asyncEnabled        bool      // Whether async rotation is enabled
}

// NewRotationManager creates a new rotation manager for the given logger
func NewRotationManager(logger *Logger, config *LoggerConfig, basePath string) *RotationManager {
	rm := &RotationManager{
		logger:            logger,
		config:            config,
		basePath:          basePath,
		sizeSyncInterval:  10 * time.Second, // Sync cached size every 10 seconds
		lastSizeSync:      time.Now(),
		asyncRotationChan: make(chan rotationRequest, 1), // Buffer of 1 to prevent blocking
		asyncEnabled:      true, // Enable async rotation by default
	}

	// Initialize cached file size
	rm.syncFileSize()

	// Initialize list of existing rotated files
	rm.scanExistingRotatedFiles()

	// Start async rotation worker
	go rm.asyncRotationWorker()

	return rm
}

// ShouldRotate checks if rotation is needed for the given entry size
func (rm *RotationManager) ShouldRotate(newEntrySize int64) bool {
	if !rm.config.RotationEnabled || rm.pendingRotation {
		return false
	}

	if rm.config.MaxFileSize <= 0 {
		return false // Unlimited file size
	}

	// Sync cached size if enough time has passed
	if time.Since(rm.lastSizeSync) > rm.sizeSyncInterval {
		rm.syncFileSize()
	}

	// Use cached size for performance
	wouldExceed := rm.cachedFileSize+newEntrySize > rm.config.MaxFileSize

	// If we're close to the limit, do a real-time check for accuracy
	if wouldExceed && rm.cachedFileSize > rm.config.MaxFileSize*8/10 {
		rm.syncFileSize() // Force sync when close to limit
		wouldExceed = rm.cachedFileSize+newEntrySize > rm.config.MaxFileSize
	}

	return wouldExceed
}

// PerformRotation rotates the current log file and creates a new one
func (rm *RotationManager) PerformRotation() error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Prevent duplicate rotations
	if rm.pendingRotation {
		return nil
	}
	rm.pendingRotation = true
	defer func() { rm.pendingRotation = false }()

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
		rm.logger.Warn("rotation_cleanup", "Failed to cleanup old files", WithError(err))
	}

	// Create new log file
	newFile, err := os.OpenFile(rm.basePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	// Update logger with new file and reset cached sizes
	rm.logger.file = newFile
	rm.logger.currentSize = 0
	rm.cachedFileSize = 0
	rm.lastSizeSync = time.Now()

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
		rm.logger.Warn("config_update_cleanup", "Failed to cleanup files after config update", WithError(err))
	}
}

// syncFileSize synchronizes the cached file size with the actual file size on disk
func (rm *RotationManager) syncFileSize() {
	if stat, err := os.Stat(rm.basePath); err == nil {
		rm.cachedFileSize = stat.Size()
		rm.logger.currentSize = stat.Size() // Keep logger's size in sync too
	}
	rm.lastSizeSync = time.Now()
}

// updateCachedSize updates the cached size after writing data
func (rm *RotationManager) updateCachedSize(deltaSize int64) {
	rm.cachedFileSize += deltaSize
}

// PerformRotationAsync performs rotation asynchronously and returns immediately
func (rm *RotationManager) PerformRotationAsync() <-chan error {
	response := make(chan error, 1)
	
	if !rm.asyncEnabled {
		// Fall back to synchronous rotation
		go func() {
			response <- rm.PerformRotation()
		}()
		return response
	}

	request := rotationRequest{
		force:    false,
		response: response,
	}

	// Try to send request (non-blocking)
	select {
	case rm.asyncRotationChan <- request:
		// Request sent successfully
	default:
		// Channel is full, fall back to sync rotation
		go func() {
			response <- rm.PerformRotation()
		}()
	}

	return response
}

// ForceRotationAsync performs forced rotation asynchronously
func (rm *RotationManager) ForceRotationAsync() <-chan error {
	response := make(chan error, 1)
	
	request := rotationRequest{
		force:    true,
		response: response,
	}

	// For forced rotation, always try async first
	select {
	case rm.asyncRotationChan <- request:
		// Request sent successfully
	default:
		// Channel is full, fall back to immediate sync rotation
		go func() {
			response <- rm.PerformRotation()
		}()
	}

	return response
}

// asyncRotationWorker handles async rotation requests
func (rm *RotationManager) asyncRotationWorker() {
	for request := range rm.asyncRotationChan {
		err := rm.PerformRotation()
		
		// Send response back
		select {
		case request.response <- err:
		default:
			// Response channel might be closed, ignore
		}
	}
}

// SetAsyncRotation enables or disables async rotation
func (rm *RotationManager) SetAsyncRotation(enabled bool) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.asyncEnabled = enabled
}

// Close shuts down the rotation manager and its background worker
func (rm *RotationManager) Close() {
	close(rm.asyncRotationChan)
}
