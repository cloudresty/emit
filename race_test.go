package emit

import (
	"strconv"
	"sync"
	"testing"
)

// TestConcurrentStructuredFields tests the race condition fix
func TestConcurrentStructuredFields(t *testing.T) {
	// Enable info level logging to ensure the logging code path is executed
	SetLevel("info")
	SetFormat("json")
	SetOutputToDiscard() // Suppress output for test

	var wg sync.WaitGroup
	numGoroutines := 100
	numLogsPerGoroutine := 100

	// Launch multiple goroutines that log concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numLogsPerGoroutine; j++ {
				Info.StructuredFields("Concurrent test message",
					ZString("goroutine_id", strconv.Itoa(goroutineID)),
					ZInt("iteration", j),
					ZString("component", "race-test"),
					ZBool("is_test", true),
				)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// If we get here without a race condition, the test passes
	t.Logf("Successfully completed %d concurrent logging operations", numGoroutines*numLogsPerGoroutine)
}

// TestConcurrentMixedLogging tests mixed logging methods under concurrent load
func TestConcurrentMixedLogging(t *testing.T) {
	SetLevel("info")
	SetFormat("json")
	SetOutputToDiscard()

	var wg sync.WaitGroup
	numGoroutines := 50

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Mix different logging methods to stress test
			Info.StructuredFields("Mixed method test 1",
				ZString("method", "structured_fields"),
				ZInt("goroutine", goroutineID),
			)

			Info.Msg("Simple message from goroutine " + strconv.Itoa(goroutineID))

			Info.StructuredFields("Mixed method test 2",
				ZString("method", "structured_fields_again"),
				ZInt("goroutine", goroutineID),
				ZBool("complex", true),
				ZFloat64("value", 3.14159),
			)
		}(i)
	}

	wg.Wait()
	t.Log("Mixed concurrent logging completed successfully")
}
