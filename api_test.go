package emit

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestCleanAPI tests the API of the emit package
func TestCleanAPI(t *testing.T) {
	var buf bytes.Buffer

	// Create a test logger
	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		component:       "test",
		version:         "1.0",
	}

	// Replace default logger temporarily
	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Test Info.Field()
	Info.Field("Test structured logging",
		NewFields().
			String("key1", "value1").
			Int("key2", 42).
			Bool("key3", true))

	// Test Error.KeyValue()
	Error.KeyValue("Test key-value logging",
		"error", "test_error",
		"code", 500)

	// Test Warn.StructuredFields()
	Warn.StructuredFields("Test structured fields logging",
		ZString("service", "test"),
		ZInt("count", 100))

	// Test Debug.Pool()
	Debug.Pool("Test pooled logging", func(pf *PooledFields) {
		pf.String("operation", "test").
			Int("duration", 123)
	})

	// Test simple message logging
	Info.Msg("Simple info message")
	Error.Msg("Simple error message")
	Warn.Msg("Simple warn message")
	Debug.Msg("Simple debug message")

	output := buf.String()

	// Verify all log entries were created
	expectedMessages := []string{
		"Test structured logging",
		"Test key-value logging",
		"Test structured fields logging",
		"Test pooled logging",
		"Simple info message",
		"Simple error message",
		"Simple warn message",
		"Simple debug message",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected message '%s' not found in output", msg)
		}
	}

	// Verify log levels
	expectedLevels := []string{"info", "error", "warn", "debug"}
	for _, level := range expectedLevels {
		if !strings.Contains(output, `"level":"`+level+`"`) {
			t.Errorf("Expected level '%s' not found in output", level)
		}
	}

	// Verify structured fields
	if !strings.Contains(output, `"key1":"value1"`) {
		t.Errorf("Expected structured field not found in output")
	}
	if !strings.Contains(output, `"key2":42`) {
		t.Errorf("Expected structured field not found in output")
	}
	if !strings.Contains(output, `"key3":true`) {
		t.Errorf("Expected structured field not found in output")
	}
}

// TestLogLevels tests log level filtering
func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           WARN, // Only WARN and ERROR should be logged
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	Debug.Msg("Debug message")  // Should be filtered out
	Info.Msg("Info message")    // Should be filtered out
	Warn.Msg("Warning message") // Should be logged
	Error.Msg("Error message")  // Should be logged

	output := buf.String()

	// Check that debug and info are NOT present
	if strings.Contains(output, "Debug message") {
		t.Errorf("Debug message should be filtered out but was found: %s", output)
	}

	if strings.Contains(output, "Info message") {
		t.Errorf("Info message should be filtered out but was found: %s", output)
	}

	// Check that warn and error ARE present
	if !strings.Contains(output, "Warning message") {
		t.Errorf("Warning message should be present but was not found: %s", output)
	}

	if !strings.Contains(output, "Error message") {
		t.Errorf("Error message should be present but was not found: %s", output)
	}
}

// TestFieldCreation tests the clean field creation functions
func TestFieldCreation(t *testing.T) {
	// Test NewField()
	fields := NewFields().
		String("test_string", "value").
		Int("test_int", 42).
		Bool("test_bool", true).
		Time("test_time", time.Now())

	fieldMap := fields.ToMap()

	if fieldMap["test_string"] != "value" {
		t.Errorf("Expected test_string to be 'value', got %v", fieldMap["test_string"])
	}

	if fieldMap["test_int"] != 42 {
		t.Errorf("Expected test_int to be 42, got %v", fieldMap["test_int"])
	}

	if fieldMap["test_bool"] != true {
		t.Errorf("Expected test_bool to be true, got %v", fieldMap["test_bool"])
	}

	// Test PooledField()
	pf := NewPooledFields()
	pf.String("pooled_test", "pooled_value")
	pooledMap := pf.ToMap()

	if pooledMap["pooled_test"] != "pooled_value" {
		t.Errorf("Expected pooled_test to be 'pooled_value', got %v", pooledMap["pooled_test"])
	}

	pf.Release() // Clean up
}

// TestParseKeyValuePairs tests the key-value parsing
func TestParseKeyValuePairs(t *testing.T) {
	// Test normal case
	result := parseKeyValuePairs("key1", "value1", "key2", 42)

	if result["key1"] != "value1" {
		t.Errorf("Expected key1 to be 'value1', got %v", result["key1"])
	}

	if result["key2"] != 42 {
		t.Errorf("Expected key2 to be 42, got %v", result["key2"])
	}

	// Test odd number of arguments (should add missing value)
	result = parseKeyValuePairs("key1", "value1", "key2")

	if result["key2"] != "<missing_value>" {
		t.Errorf("Expected key2 to have missing_value placeholder, got %v", result["key2"])
	}

	// Test empty case
	result = parseKeyValuePairs()
	if result != nil {
		t.Errorf("Expected nil result for empty args, got %v", result)
	}
}
