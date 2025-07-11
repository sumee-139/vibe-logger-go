package main

import (
	"fmt"
	"time"
	
	"github.com/fladdict/vibe-logger-go/pkg/vibelogger"
)

func main() {
	// Create a file logger
	logger, err := vibelogger.CreateFileLogger("example_app")
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer logger.Close()
	
	// Basic logging
	logger.Info("app_startup", "Application started successfully")
	
	// Logging with context
	logger.Info("user_login", "User logged in",
		vibelogger.WithContext(map[string]interface{}{
			"user_id": "user123",
			"ip_address": "192.168.1.100",
			"user_agent": "Mozilla/5.0...",
		}),
	)
	
	// Logging with AI instructions
	logger.Info("data_processing", "Processing user data",
		vibelogger.WithContext(map[string]interface{}{
			"record_count": 1500,
			"processing_time": "2.3s",
		}),
		vibelogger.WithHumanNote("This process might be slow with large datasets"),
		vibelogger.WithAITodo("AI-TODO: Consider implementing batch processing for better performance"),
	)
	
	// Error logging with correlation ID
	correlationID := "req-789-xyz"
	logger.Error("database_connection", "Failed to connect to database",
		vibelogger.WithCorrelationID(correlationID),
		vibelogger.WithContext(map[string]interface{}{
			"database_host": "localhost:5432",
			"retry_count": 3,
		}),
		vibelogger.WithHumanNote("Database connection issues started after the recent update"),
		vibelogger.WithAITodo("AI-TODO: Check database connection pool settings and recent configuration changes"),
	)
	
	// Logging with duration
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // Simulate some work
	duration := time.Since(start)
	
	logger.Info("api_call", "External API call completed",
		vibelogger.WithDuration(duration),
		vibelogger.WithContext(map[string]interface{}{
			"api_endpoint": "https://api.example.com/users",
			"status_code": 200,
			"response_size": 1024,
		}),
	)
	
	// Warning with user context
	logger.Warn("rate_limit", "API rate limit approaching",
		vibelogger.WithUserID("user456"),
		vibelogger.WithContext(map[string]interface{}{
			"current_requests": 950,
			"limit": 1000,
			"reset_time": time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}),
		vibelogger.WithAITodo("AI-TODO: Implement rate limiting warnings and suggest alternative endpoints"),
	)
	
	// Debug logging
	logger.Debug("cache_operation", "Cache miss occurred",
		vibelogger.WithContext(map[string]interface{}{
			"cache_key": "user:profile:123",
			"cache_type": "redis",
			"miss_reason": "key_expired",
		}),
	)
	
	fmt.Println("Example completed. Check the logs/ directory for the generated log file.")
}