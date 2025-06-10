package emit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestFields tests the Fields builder pattern
func TestFields(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           INFO,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Test Fields builder pattern
	fields := F().
		String("username", "john_doe").
		Int("user_id", 12345).
		String("email", "john@example.com").
		String("password", "secret123").
		Bool("active", true).
		Float64("score", 95.5)

	InfoF("User action", fields)

	output := buf.String()

	// Parse JSON to check fields
	var logEntry map[string]any
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		err := json.Unmarshal([]byte(lines[0]), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		logFields, ok := logEntry["fields"].(map[string]any)
		if !ok {
			t.Fatalf("Fields not found in log entry")
		}

		// Check that fields are present and masked correctly
		if logFields["username"] != "***PII***" {
			t.Errorf("Expected username to be masked as PII, got: %v", logFields["username"])
		}

		if logFields["user_id"] != float64(12345) { // JSON unmarshals numbers as float64
			t.Errorf("Expected user_id 12345, got: %v", logFields["user_id"])
		}

		if logFields["email"] != "***PII***" {
			t.Errorf("Expected email to be masked, got: %v", logFields["email"])
		}

		if logFields["password"] != "***MASKED***" {
			t.Errorf("Expected password to be masked, got: %v", logFields["password"])
		}

		if logFields["active"] != true {
			t.Errorf("Expected active true, got: %v", logFields["active"])
		}

		if logFields["score"] != 95.5 {
			t.Errorf("Expected score 95.5, got: %v", logFields["score"])
		}
	}
}

// TestFieldsChaining tests method chaining
func TestFieldsChaining(t *testing.T) {
	// Test that all methods return Fields for chaining
	fields := NewFields().
		Set("key1", "value1").
		Add("key2", "value2").
		With("key3", "value3").
		String("key4", "value4").
		Int("key5", 5).
		Bool("key6", true)

	if len(fields) != 6 {
		t.Errorf("Expected 6 fields, got %d", len(fields))
	}

	expected := map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": 5,
		"key6": true,
	}

	for key, expectedValue := range expected {
		if fields[key] != expectedValue {
			t.Errorf("Expected %s=%v, got %v", key, expectedValue, fields[key])
		}
	}
}

// TestVariadicKeyValue tests the key-value pair API
func TestVariadicKeyValue(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           INFO,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Test variadic key-value API
	InfoKV("User login",
		"username", "jane_doe",
		"user_id", 67890,
		"email", "jane@example.com",
		"password", "secret456",
		"timestamp", "2025-06-10T09:00:00Z",
	)

	output := buf.String()

	// Parse JSON to check fields
	var logEntry map[string]any
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		err := json.Unmarshal([]byte(lines[0]), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		logFields, ok := logEntry["fields"].(map[string]any)
		if !ok {
			t.Fatalf("Fields not found in log entry")
		}

		// Check that fields are present and masked correctly
		if logFields["username"] != "***PII***" {
			t.Errorf("Expected username to be masked as PII, got: %v", logFields["username"])
		}

		if logFields["user_id"] != float64(67890) {
			t.Errorf("Expected user_id 67890, got: %v", logFields["user_id"])
		}

		if logFields["email"] != "***PII***" {
			t.Errorf("Expected email to be masked, got: %v", logFields["email"])
		}

		if logFields["password"] != "***MASKED***" {
			t.Errorf("Expected password to be masked, got: %v", logFields["password"])
		}

		if logFields["timestamp"] != "2025-06-10T09:00:00Z" {
			t.Errorf("Expected timestamp 2025-06-10T09:00:00Z, got: %v", logFields["timestamp"])
		}
	}
}

// TestQuickFieldHelpers tests quick field creation helpers
func TestQuickFieldHelpers(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           INFO,
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

	// Test quick field helpers
	InfoF("Single field test", StringField("name", "test"))

	buf.Reset()
	InfoF("Int field test", IntField("count", 42))

	buf.Reset()
	InfoF("Multiple fields test", Field("key1", "value1").Set("key2", 123))

	output := buf.String()
	if !strings.Contains(output, `"key1":"value1"`) {
		t.Errorf("Expected key1:value1 in output: %s", output)
	}
	if !strings.Contains(output, `"key2":123`) {
		t.Errorf("Expected key2:123 in output: %s", output)
	}
}

// TestFieldsMerge tests merging Fields objects
func TestFieldsMerge(t *testing.T) {
	fields1 := F().String("name", "John").Int("age", 30)
	fields2 := F().String("city", "NYC").String("country", "USA")

	merged := fields1.Merge(fields2)

	if len(merged) != 4 {
		t.Errorf("Expected 4 fields after merge, got %d", len(merged))
	}

	expected := map[string]any{
		"name":    "John",
		"age":     30,
		"city":    "NYC",
		"country": "USA",
	}

	for key, expectedValue := range expected {
		if merged[key] != expectedValue {
			t.Errorf("Expected %s=%v, got %v", key, expectedValue, merged[key])
		}
	}
}

// TestFieldsClone tests cloning Fields objects
func TestFieldsClone(t *testing.T) {
	original := F().String("name", "John").Int("age", 30)
	clone := original.Clone()

	// Modify original
	original.String("city", "NYC")

	// Clone should not be affected
	if len(clone) != 2 {
		t.Errorf("Expected clone to have 2 fields, got %d", len(clone))
	}

	if len(original) != 3 {
		t.Errorf("Expected original to have 3 fields, got %d", len(original))
	}

	if _, exists := clone["city"]; exists {
		t.Errorf("Clone should not have city field")
	}
}

// TestOddNumberOfKeyValuePairs tests handling of odd number of arguments
func TestOddNumberOfKeyValuePairs(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           INFO,
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

	// Test with odd number of arguments (should handle gracefully)
	InfoKV("Test odd args", "key1", "value1", "key2") // Missing value for key2

	output := buf.String()

	// Parse JSON to check fields
	var logEntry map[string]any
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		err := json.Unmarshal([]byte(lines[0]), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		logFields, ok := logEntry["fields"].(map[string]any)
		if !ok {
			t.Fatalf("Fields not found in log entry")
		}

		// Should have key1 with value1 and key2 with placeholder
		if logFields["key1"] != "value1" {
			t.Errorf("Expected key1:value1, got: %v", logFields["key1"])
		}

		if logFields["key2"] != "<missing_value>" {
			t.Errorf("Expected key2:<missing_value>, got: %v", logFields["key2"])
		}
	}
}

// TestAllLogLevelsWithFields tests all log levels with Fields
func TestAllLogLevelsWithFields(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
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

	fields := F().String("test", "value")

	// Test all log levels with Fields
	DebugF("Debug message", fields)
	InfoF("Info message", fields)
	WarnF("Warn message", fields)
	ErrorF("Error message", fields)

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 4 {
		t.Errorf("Expected 4 log lines, got %d", len(lines))
	}

	expectedLevels := []string{"debug", "info", "warn", "error"}
	for i, line := range lines {
		if !strings.Contains(line, `"level":"`+expectedLevels[i]+`"`) {
			t.Errorf("Expected level %s in line %d: %s", expectedLevels[i], i, line)
		}
		if !strings.Contains(line, `"test":"value"`) {
			t.Errorf("Expected field test:value in line %d: %s", i, line)
		}
	}
}

// TestAllLogLevelsWithKeyValue tests all log levels with key-value pairs
func TestAllLogLevelsWithKeyValue(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
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

	// Test all log levels with key-value pairs
	DebugKV("Debug message", "test", "value")
	InfoKV("Info message", "test", "value")
	WarnKV("Warn message", "test", "value")
	ErrorKV("Error message", "test", "value")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 4 {
		t.Errorf("Expected 4 log lines, got %d", len(lines))
	}

	expectedLevels := []string{"debug", "info", "warn", "error"}
	for i, line := range lines {
		if !strings.Contains(line, `"level":"`+expectedLevels[i]+`"`) {
			t.Errorf("Expected level %s in line %d: %s", expectedLevels[i], i, line)
		}
		if !strings.Contains(line, `"test":"value"`) {
			t.Errorf("Expected field test:value in line %d: %s", i, line)
		}
	}
}
