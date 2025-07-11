package vibelogger

import (
	"fmt"
	"time"
)

// LogOption is a function that modifies a LogEntry
type LogOption func(*LogEntry)

// WithContext adds context information to the log entry
func WithContext(context map[string]interface{}) LogOption {
	return func(entry *LogEntry) {
		if entry.Context == nil {
			entry.Context = make(map[string]interface{})
		}
		for k, v := range context {
			entry.Context[k] = v
		}
	}
}

// WithHumanNote adds a human-readable note for AI analysis
func WithHumanNote(note string) LogOption {
	return func(entry *LogEntry) {
		entry.HumanNote = note
	}
}

// WithAITodo adds an AI todo instruction
func WithAITodo(todo string) LogOption {
	return func(entry *LogEntry) {
		entry.AITodo = todo
	}
}

// WithCorrelationID adds a correlation ID for tracking related logs
func WithCorrelationID(id string) LogOption {
	return func(entry *LogEntry) {
		entry.CorrelationID = id
	}
}

// WithFields is a convenience function for adding multiple context fields
func WithFields(fields map[string]interface{}) LogOption {
	return WithContext(fields)
}

// WithError adds error information to the context
func WithError(err error) LogOption {
	return func(entry *LogEntry) {
		if entry.Context == nil {
			entry.Context = make(map[string]interface{})
		}
		entry.Context["error"] = err.Error()
		entry.Context["error_type"] = fmt.Sprintf("%T", err)
	}
}

// WithUserID adds user ID to the context
func WithUserID(userID string) LogOption {
	return func(entry *LogEntry) {
		if entry.Context == nil {
			entry.Context = make(map[string]interface{})
		}
		entry.Context["user_id"] = userID
	}
}

// WithRequestID adds request ID to the context
func WithRequestID(requestID string) LogOption {
	return func(entry *LogEntry) {
		if entry.Context == nil {
			entry.Context = make(map[string]interface{})
		}
		entry.Context["request_id"] = requestID
	}
}

// WithDuration adds duration information to the context
func WithDuration(duration time.Duration) LogOption {
	return func(entry *LogEntry) {
		if entry.Context == nil {
			entry.Context = make(map[string]interface{})
		}
		entry.Context["duration_ms"] = duration.Milliseconds()
		entry.Context["duration_human"] = duration.String()
	}
}
