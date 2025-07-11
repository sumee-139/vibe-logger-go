package vibelogger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	Timestamp    time.Time              `json:"timestamp"`
	Level        LogLevel               `json:"level"`
	Operation    string                 `json:"operation"`
	Message      string                 `json:"message"`
	Context      map[string]interface{} `json:"context,omitempty"`
	HumanNote    string                 `json:"human_note,omitempty"`
	AITodo       string                 `json:"ai_todo,omitempty"`
	StackTrace   []string               `json:"stack_trace,omitempty"`
	Environment  map[string]string      `json:"environment,omitempty"`
	CorrelationID string                `json:"correlation_id,omitempty"`
}

// Logger is the main vibe logger instance
type Logger struct {
	name     string
	filePath string
	file     *os.File
	mutex    sync.Mutex
}

// NewLogger creates a new Logger instance
func NewLogger(name string) *Logger {
	return &Logger{
		name: name,
	}
}

// CreateFileLogger creates a new file-based logger
func CreateFileLogger(name string) (*Logger, error) {
	logger := NewLogger(name)
	
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}
	
	// Create timestamped log file
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.log", name, timestamp)
	logger.filePath = filepath.Join(logDir, filename)
	
	file, err := os.OpenFile(logger.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}
	
	logger.file = file
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
	
	if l.file != nil {
		if _, err := l.file.Write(jsonData); err != nil {
			return fmt.Errorf("failed to write to log file: %w", err)
		}
		if _, err := l.file.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline to log file: %w", err)
		}
	}
	
	// Also output to console
	fmt.Printf("%s\n", string(jsonData))
	
	return nil
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