package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sumee-139/vibe-logger-go/pkg/vibelogger"
)

func main() {
	// Create a file logger
	logger, err := vibelogger.CreateFileLogger("ai_test_demo")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	fmt.Println("=== AI-Optimized Logger Demo ===")
	fmt.Println()

	// Test various AI-optimized patterns
	testCases := []struct {
		level     string
		operation string
		message   string
		options   []vibelogger.LogOption
	}{
		// Database errors
		{"error", "db_query", "connection refused by database server", nil},
		{"error", "user_lookup", "no rows found for user ID 123", nil},

		// Network errors
		{"warn", "api_call", "timeout occurred after 5 seconds", nil},
		{"error", "http_request", "received 502 bad gateway", nil},

		// Authentication errors
		{"error", "user_auth", "unauthorized access attempt", nil},
		{"warn", "token_check", "JWT token expired", nil},

		// File system errors
		{"error", "file_read", "file not found: /etc/config.json", nil},
		{"error", "log_write", "permission denied to write log file", nil},

		// Performance issues
		{"warn", "query_execution", "slow query detected - 2.5 seconds", nil},
		{"error", "memory_check", "high memory usage detected", nil},

		// Validation errors
		{"error", "user_input", "invalid email format provided", nil},
		{"warn", "data_validation", "malformed JSON in request body", nil},

		// Business logic with context
		{"info", "user_registration", "new user account created", []vibelogger.LogOption{
			vibelogger.WithContext(map[string]interface{}{
				"user_id": "12345",
				"email":   "test@example.com",
				"source":  "web_app",
			}),
			vibelogger.WithHumanNote("User registration successful - verify email sent"),
			vibelogger.WithCorrelationID("reg-2025-001"),
		}},

		// AI Todo example
		{"warn", "payment_processing", "payment gateway response time increasing", []vibelogger.LogOption{
			vibelogger.WithAITodo("Monitor payment gateway performance and consider fallback options"),
			vibelogger.WithDuration(2500 * time.Millisecond),
		}},
	}

	for i, tc := range testCases {
		fmt.Printf("--- Test Case %d: %s ---\n", i+1, tc.operation)

		var err error
		switch tc.level {
		case "info":
			err = logger.Info(tc.operation, tc.message, tc.options...)
		case "warn":
			err = logger.Warn(tc.operation, tc.message, tc.options...)
		case "error":
			err = logger.Error(tc.operation, tc.message, tc.options...)
		case "debug":
			err = logger.Debug(tc.operation, tc.message, tc.options...)
		}

		if err != nil {
			fmt.Printf("Error logging: %v\n", err)
		}

		fmt.Println()
		time.Sleep(100 * time.Millisecond) // Small delay for readability
	}

	fmt.Println("=== Demo Complete ===")
	fmt.Printf("Check the log file in the 'logs' directory for complete AI-optimized output\n")
}
