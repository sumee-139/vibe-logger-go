package vibelogger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	DEBUG LogLevel = "DEBUG"
)

// LogEntry represents a single log entry with AI-optimized structure
type LogEntry struct {
	Timestamp     time.Time              `json:"timestamp"`
	Level         LogLevel               `json:"level"`
	Operation     string                 `json:"operation"`
	Message       string                 `json:"message"`
	Context       map[string]interface{} `json:"context,omitempty"`
	HumanNote     string                 `json:"human_note,omitempty"`
	AITodo        string                 `json:"ai_todo,omitempty"`
	StackTrace    []string               `json:"stack_trace,omitempty"`
	Environment   map[string]string      `json:"environment,omitempty"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	// AI-optimized fields
	Severity   int    `json:"severity"`             // 1-5 scale for AI prioritization
	Category   string `json:"category,omitempty"`   // business_logic, system, user_action, etc.
	Searchable string `json:"searchable,omitempty"` // AI-friendly search terms
	Pattern    string `json:"pattern,omitempty"`    // Known error patterns
	Suggestion string `json:"suggestion,omitempty"` // AI debugging suggestions
}

// Logger is the main vibe logger instance
type Logger struct {
	name        string
	filePath    string
	file        *os.File
	mutex       sync.Mutex
	config      *LoggerConfig
	currentSize int64
	memoryLogs  []LogEntry
	memoryMutex sync.Mutex
	rotationMgr *RotationManager
}

// NewLogger creates a new Logger instance with default configuration
func NewLogger(name string) *Logger {
	return &Logger{
		name:   name,
		config: DefaultConfig(),
	}
}

// NewLoggerWithConfig creates a new Logger instance with custom configuration
func NewLoggerWithConfig(name string, config *LoggerConfig) *Logger {
	if config == nil {
		config = DefaultConfig()
	}
	config.Validate()

	return &Logger{
		name:   name,
		config: config,
	}
}

// CreateFileLogger creates a new file-based logger with default configuration
func CreateFileLogger(name string) (*Logger, error) {
	return CreateFileLoggerWithConfig(name, DefaultConfig())
}

// CreateFileLoggerWithConfig creates a new file-based logger with custom configuration
func CreateFileLoggerWithConfig(name string, config *LoggerConfig) (*Logger, error) {
	logger := NewLoggerWithConfig(name, config)

	// Use custom file path or generate default
	var logDir, filename string
	if config.FilePath != "" {
		logger.filePath = config.FilePath
		// Create directory for custom file path if it doesn't exist
		dir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for custom file path: %w", err)
		}
	} else {
		// Create logs directory if it doesn't exist
		logDir = "logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create logs directory: %w", err)
		}

		// Create timestamped log file
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("%s_%s.log", name, timestamp)
		logger.filePath = filepath.Join(logDir, filename)
	}

	// Open or create the log file
	file, err := os.OpenFile(logger.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	// Get current file size for MaxFileSize tracking
	if stat, err := file.Stat(); err == nil {
		logger.currentSize = stat.Size()
	}

	logger.file = file

	// Initialize rotation manager if rotation is enabled
	if config.RotationEnabled {
		logger.rotationMgr = NewRotationManager(logger, config, logger.filePath)
	}

	return logger, nil
}

// Log writes a log entry with the specified level
func (l *Logger) Log(level LogLevel, operation, message string, options ...LogOption) error {
	entry := LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Operation: operation,
		Message:   message,
		Context:   make(map[string]interface{}),
	}

	// Apply options
	for _, opt := range options {
		opt(&entry)
	}

	// Add stack trace for ERROR level
	if level == ERROR {
		entry.StackTrace = getStackTrace()
	}

	// Add environment information
	entry.Environment = getEnvironment()

	// Set AI-optimized fields
	entry.Severity = getSeverityScore(level)
	entry.Category = inferCategory(operation, message)
	entry.Searchable = generateSearchableTerms(operation, message)
	entry.Pattern = detectKnownPattern(operation, message)
	entry.Suggestion = generateAISuggestion(level, operation, message)

	return l.writeEntry(entry)
}

// Info logs an info level message
func (l *Logger) Info(operation, message string, options ...LogOption) error {
	return l.Log(INFO, operation, message, options...)
}

// Warn logs a warning level message
func (l *Logger) Warn(operation, message string, options ...LogOption) error {
	return l.Log(WARN, operation, message, options...)
}

// Error logs an error level message
func (l *Logger) Error(operation, message string, options ...LogOption) error {
	return l.Log(ERROR, operation, message, options...)
}

// Debug logs a debug level message
func (l *Logger) Debug(operation, message string, options ...LogOption) error {
	return l.Log(DEBUG, operation, message, options...)
}

// Close closes the logger and its file handle
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file != nil {
		err := l.file.Close()
		l.file = nil // Set to nil to prevent double-close
		return err
	}
	return nil
}

// writeEntry writes a log entry to the file
func (l *Logger) writeEntry(entry LogEntry) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Add to memory log if enabled
	if l.config.EnableMemoryLog {
		l.addToMemoryLog(entry)
	}

	// Write to file if AutoSave is enabled and file exists
	if l.config.AutoSave && l.file != nil {
		entrySize := int64(len(jsonData) + 1) // +1 for newline

		// Check if rotation is needed and perform it
		if l.rotationMgr != nil && l.rotationMgr.ShouldRotate(entrySize) {
			if err := l.rotationMgr.PerformRotation(); err != nil {
				return fmt.Errorf("failed to rotate log file: %w", err)
			}
		}

		if _, err := l.file.Write(jsonData); err != nil {
			return fmt.Errorf("failed to write to log file: %w", err)
		}
		if _, err := l.file.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline to log file: %w", err)
		}

		// Update current file size
		l.currentSize += entrySize
	}

	// Always output to console for debugging
	fmt.Printf("%s\n", string(jsonData))

	return nil
}

// addToMemoryLog adds an entry to the in-memory log
func (l *Logger) addToMemoryLog(entry LogEntry) {
	l.memoryMutex.Lock()
	defer l.memoryMutex.Unlock()

	l.memoryLogs = append(l.memoryLogs, entry)

	// Enforce memory log limit
	if l.config.MemoryLogLimit > 0 && len(l.memoryLogs) > l.config.MemoryLogLimit {
		// Remove oldest entries
		excess := len(l.memoryLogs) - l.config.MemoryLogLimit
		l.memoryLogs = l.memoryLogs[excess:]
	}
}

// GetMemoryLogs returns a copy of the current memory logs
func (l *Logger) GetMemoryLogs() []LogEntry {
	l.memoryMutex.Lock()
	defer l.memoryMutex.Unlock()

	// Return a copy to prevent external modification
	logs := make([]LogEntry, len(l.memoryLogs))
	copy(logs, l.memoryLogs)
	return logs
}

// ClearMemoryLogs clears all entries from the memory log
func (l *Logger) ClearMemoryLogs() {
	l.memoryMutex.Lock()
	defer l.memoryMutex.Unlock()
	l.memoryLogs = nil
}

// getStackTrace returns the current stack trace
func getStackTrace() []string {
	var stack []string

	// Skip the first 3 frames (getStackTrace, writeEntry, Log)
	for i := 3; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}

		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}

	return stack
}

// getEnvironment returns current environment information
func getEnvironment() map[string]string {
	return map[string]string{
		"go_version": runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"pid":        fmt.Sprintf("%d", os.Getpid()),
		"pwd":        func() string { pwd, _ := os.Getwd(); return pwd }(),
	}
}

// ForceRotation manually triggers log file rotation
func (l *Logger) ForceRotation() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.rotationMgr == nil {
		return fmt.Errorf("rotation is not enabled")
	}

	return l.rotationMgr.PerformRotation()
}

// GetRotatedFiles returns the list of current rotated files
func (l *Logger) GetRotatedFiles() []string {
	if l.rotationMgr == nil {
		return nil
	}

	return l.rotationMgr.GetRotatedFiles()
}

// UpdateConfig updates the logger configuration including rotation settings
func (l *Logger) UpdateConfig(config *LoggerConfig) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Validate new configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	l.config = config

	// Initialize or update rotation manager
	if config.RotationEnabled && l.rotationMgr == nil {
		// Enable rotation
		l.rotationMgr = NewRotationManager(l, config, l.filePath)
	} else if !config.RotationEnabled && l.rotationMgr != nil {
		// Disable rotation
		l.rotationMgr = nil
	} else if l.rotationMgr != nil {
		// Update existing rotation manager
		l.rotationMgr.UpdateConfig(config)
	}

	return nil
}

// getSeverityScore converts log level to numerical severity for AI prioritization
func getSeverityScore(level LogLevel) int {
	switch level {
	case DEBUG:
		return 1
	case INFO:
		return 2
	case WARN:
		return 3
	case ERROR:
		return 4
	default:
		return 2
	}
}

// inferCategory tries to determine the category of the log entry based on operation and message
func inferCategory(operation, message string) string {
	operation = fmt.Sprintf("%s %s", operation, message)

	// Simple keyword-based categorization
	if containsAny(operation, []string{"user", "login", "auth", "session"}) {
		return "user_action"
	}
	if containsAny(operation, []string{"db", "database", "sql", "query"}) {
		return "database"
	}
	if containsAny(operation, []string{"api", "http", "request", "response"}) {
		return "api"
	}
	if containsAny(operation, []string{"file", "disk", "io", "read", "write"}) {
		return "system"
	}
	if containsAny(operation, []string{"business", "logic", "validation", "calculation"}) {
		return "business_logic"
	}

	return "general"
}

// generateSearchableTerms creates AI-friendly search terms from operation and message
func generateSearchableTerms(operation, message string) string {
	combined := fmt.Sprintf("%s %s", operation, message)

	// Extract key terms for AI searchability
	keywords := []string{}

	// Add operation as a keyword
	if operation != "" {
		keywords = append(keywords, operation)
	}

	// Add category-specific keywords
	if containsAny(combined, []string{"error", "failed", "exception", "panic"}) {
		keywords = append(keywords, "error")
	}
	if containsAny(combined, []string{"start", "begin", "init", "startup"}) {
		keywords = append(keywords, "start")
	}
	if containsAny(combined, []string{"end", "complete", "finish", "done"}) {
		keywords = append(keywords, "complete")
	}
	if containsAny(combined, []string{"timeout", "slow", "performance"}) {
		keywords = append(keywords, "performance")
	}
	if containsAny(combined, []string{"retry", "attempt", "fallback"}) {
		keywords = append(keywords, "retry")
	}

	return strings.Join(append(keywords, combined), " ")
}

// detectKnownPattern identifies common error patterns that AI can recognize
func detectKnownPattern(operation, message string) string {
	combined := strings.ToLower(fmt.Sprintf("%s %s", operation, message))

	// Database patterns
	if containsAny(combined, []string{"connection refused", "connection timeout", "no rows", "duplicate key"}) {
		return "database_error"
	}

	// Network patterns
	if containsAny(combined, []string{"network unreachable", "connection reset", "timeout", "502", "503", "504"}) {
		return "network_error"
	}

	// Authentication patterns
	if containsAny(combined, []string{"unauthorized", "forbidden", "invalid token", "expired"}) {
		return "auth_error"
	}

	// File system patterns
	if containsAny(combined, []string{"file not found", "permission denied", "disk full", "no space"}) {
		return "filesystem_error"
	}

	// Performance patterns
	if containsAny(combined, []string{"slow query", "high memory", "cpu usage", "memory leak"}) {
		return "performance_issue"
	}

	// Validation patterns
	if containsAny(combined, []string{"invalid input", "validation failed", "bad request", "malformed", "invalid format"}) {
		return "validation_error"
	}

	return "unknown_pattern"
}

// generateAISuggestion provides AI-friendly debugging suggestions
func generateAISuggestion(level LogLevel, operation, message string) string {
	combined := strings.ToLower(fmt.Sprintf("%s %s", operation, message))

	// Only provide suggestions for warnings and errors
	if level != WARN && level != ERROR {
		return ""
	}

	// Database suggestions
	if containsAny(combined, []string{"connection refused", "connection timeout"}) {
		return "Check database connectivity and connection pool settings"
	}
	if containsAny(combined, []string{"no rows", "not found"}) {
		return "Verify query parameters and data existence"
	}

	// Network suggestions
	if containsAny(combined, []string{"timeout", "network unreachable"}) {
		return "Check network connectivity and service availability"
	}
	if containsAny(combined, []string{"502", "503", "504"}) {
		return "Verify upstream service health and load balancing"
	}

	// Authentication suggestions
	if containsAny(combined, []string{"unauthorized", "forbidden"}) {
		return "Verify authentication credentials and permissions"
	}
	if containsAny(combined, []string{"expired", "invalid token"}) {
		return "Check token expiration and refresh mechanism"
	}

	// File system suggestions
	if containsAny(combined, []string{"file not found", "no such file"}) {
		return "Verify file path and existence"
	}
	if containsAny(combined, []string{"permission denied"}) {
		return "Check file permissions and user access rights"
	}

	// Performance suggestions
	if containsAny(combined, []string{"slow", "timeout", "performance"}) {
		return "Consider optimizing query/operation or adding timeout handling"
	}

	// Validation suggestions
	if containsAny(combined, []string{"invalid", "validation", "malformed"}) {
		return "Verify input data format and validation rules"
	}

	return "Review logs and investigate root cause"
}

// containsAny checks if any of the substrings exist in the main string (case-insensitive)
func containsAny(s string, substrings []string) bool {
	s = strings.ToLower(s)
	for _, substring := range substrings {
		if strings.Contains(s, strings.ToLower(substring)) {
			return true
		}
	}
	return false
}
