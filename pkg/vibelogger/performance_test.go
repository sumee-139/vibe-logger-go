package vibelogger

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkLogger_Info(b *testing.B) {
	logger, err := CreateFileLogger("bench_info")
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark_operation", fmt.Sprintf("benchmark message %d", i))
	}
}

func BenchmarkLogger_Error(b *testing.B) {
	logger, err := CreateFileLogger("bench_error")
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("benchmark_operation", fmt.Sprintf("benchmark error %d", i))
	}
}

func BenchmarkLogger_WithContext(b *testing.B) {
	logger, err := CreateFileLogger("bench_context")
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark_operation", fmt.Sprintf("benchmark message %d", i),
			WithContext(map[string]interface{}{
				"iteration": i,
				"test":      "benchmark",
			}),
		)
	}
}

func BenchmarkLogger_Concurrent(b *testing.B) {
	logger, err := CreateFileLogger("bench_concurrent")
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Info("concurrent_operation", fmt.Sprintf("concurrent message %d", i))
			i++
		}
	})
}

func BenchmarkAIOptimization(b *testing.B) {
	logger, err := CreateFileLogger("bench_ai")
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	defer cleanup()

	testMessages := []string{
		"database connection failed",
		"network timeout occurred",
		"user authentication failed",
		"file not found error",
		"slow query detected",
		"validation error occurred",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msg := testMessages[i%len(testMessages)]
		logger.Error("ai_benchmark", msg)
	}
}

func TestConcurrentLogging(t *testing.T) {
	logger, err := CreateFileLogger("test_concurrent")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	defer cleanup()

	const numGoroutines = 10
	const messagesPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				err := logger.Info("concurrent_test",
					fmt.Sprintf("Message from goroutine %d, iteration %d", goroutineID, j),
					WithContext(map[string]interface{}{
						"goroutine_id": goroutineID,
						"iteration":    j,
					}),
				)
				if err != nil {
					t.Errorf("Failed to log: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalMessages := numGoroutines * messagesPerGoroutine
	messagesPerSecond := float64(totalMessages) / elapsed.Seconds()

	t.Logf("Concurrent test completed:")
	t.Logf("  Total messages: %d", totalMessages)
	t.Logf("  Total time: %v", elapsed)
	t.Logf("  Messages per second: %.2f", messagesPerSecond)

	// Performance threshold check
	if messagesPerSecond < 1000 {
		t.Logf("Warning: Performance below 1000 messages/sec (%.2f)", messagesPerSecond)
	}
}

func TestMemoryUsage(t *testing.T) {
	logger, err := CreateFileLogger("test_memory")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()
	defer cleanup()

	const numMessages = 1000

	// Log many messages to test memory usage
	for i := 0; i < numMessages; i++ {
		err := logger.Info("memory_test",
			fmt.Sprintf("Test message %d with some additional content to simulate real usage", i),
			WithContext(map[string]interface{}{
				"iteration": i,
				"data":      fmt.Sprintf("some_data_%d", i),
				"timestamp": time.Now().Unix(),
			}),
		)
		if err != nil {
			t.Errorf("Failed to log message %d: %v", i, err)
		}
	}

	t.Logf("Memory test completed with %d messages", numMessages)
}

func TestPatternDetectionPerformance(t *testing.T) {
	testCases := []struct {
		operation string
		message   string
		expected  string
	}{
		{"db_query", "connection refused", "database_error"},
		{"api_call", "timeout occurred", "network_error"},
		{"user_auth", "unauthorized access", "auth_error"},
		{"file_read", "file not found", "filesystem_error"},
		{"query_exec", "slow query detected", "performance_issue"},
		{"input_validation", "invalid format", "validation_error"},
		{"general_op", "some random message", "unknown_pattern"},
	}

	start := time.Now()
	const iterations = 10000

	for i := 0; i < iterations; i++ {
		for _, tc := range testCases {
			pattern := detectKnownPattern(tc.operation, tc.message)
			if pattern != tc.expected {
				t.Errorf("Pattern detection failed: expected %s, got %s", tc.expected, pattern)
			}
		}
	}

	elapsed := time.Since(start)
	operationsPerSecond := float64(len(testCases)*iterations) / elapsed.Seconds()

	t.Logf("Pattern detection performance:")
	t.Logf("  Total operations: %d", len(testCases)*iterations)
	t.Logf("  Total time: %v", elapsed)
	t.Logf("  Operations per second: %.2f", operationsPerSecond)

	// Performance threshold check
	if operationsPerSecond < 100000 {
		t.Logf("Warning: Pattern detection performance below 100k ops/sec (%.2f)", operationsPerSecond)
	}
}

func cleanup() {
	// Clean up test files
	// Note: In a real implementation, you might want to clean up test log files
	// For now, we'll let them persist for inspection
}
