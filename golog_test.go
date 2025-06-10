package emit

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// TestBasicLogging tests basic logging functionality
func TestBasicLogging(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a test logger
	testLogger := &Logger{
		level:           INFO,
		writer:          &buf,
		showCaller:      false,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
		component:       "test-app",
		version:         "v1.0.0",
	}

	// Temporarily replace default logger
	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Test basic logging
	Info("Test info message")

	output := buf.String()
	if !strings.Contains(output, "Test info message") {
		t.Errorf("Expected log message not found in output: %s", output)
	}

	if !strings.Contains(output, `"level":"info"`) {
		t.Errorf("Expected info level not found in output: %s", output)
	}

	if !strings.Contains(output, `"component":"test-app"`) {
		t.Errorf("Expected component not found in output: %s", output)
	}
}

// TestLogLevels tests different log levels
func TestLogLevels(t *testing.T) {
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

	// Test all log levels
	Debug("Debug message")
	Info("Info message")
	Warning("Warning message")
	Error("Error message")

	output := buf.String()

	expectedLevels := []string{"debug", "info", "warn", "error"}
	for _, level := range expectedLevels {
		if !strings.Contains(output, `"level":"`+level+`"`) {
			t.Errorf("Expected level %s not found in output: %s", level, output)
		}
	}
}

// TestLogLevelFiltering tests that log level filtering works
func TestLogLevelFiltering(t *testing.T) {
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

	Debug("Debug message")     // Should be filtered out
	Info("Info message")       // Should be filtered out
	Warning("Warning message") // Should be logged
	Error("Error message")     // Should be logged

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

// TestSensitiveDataMasking tests that sensitive data is properly masked
func TestSensitiveDataMasking(t *testing.T) {
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

	// Test with sensitive and PII data
	InfoWithFields("User login", map[string]any{
		"email":    "user@example.com",
		"password": "secret123",
		"user_id":  12345,
		"session":  "session-token-123",
	})

	output := buf.String()

	// Parse JSON to check fields
	var logEntry map[string]any
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		err := json.Unmarshal([]byte(lines[0]), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		fields, ok := logEntry["fields"].(map[string]any)
		if !ok {
			t.Fatalf("Fields not found in log entry")
		}

		// Check that email is masked as PII
		if fields["email"] != "***PII***" {
			t.Errorf("Expected email to be masked as ***PII***, got: %v", fields["email"])
		}

		// Check that password is masked as sensitive
		if fields["password"] != "***MASKED***" {
			t.Errorf("Expected password to be masked as ***MASKED***, got: %v", fields["password"])
		}

		// Check that session is masked as sensitive
		if fields["session"] != "***MASKED***" {
			t.Errorf("Expected session to be masked as ***MASKED***, got: %v", fields["session"])
		}

		// Check that user_id is NOT masked
		if fields["user_id"] != float64(12345) { // JSON unmarshals numbers as float64
			t.Errorf("Expected user_id to not be masked, got: %v", fields["user_id"])
		}
	}
}

// TestDataMaskingDisabled tests that masking can be disabled
func TestDataMaskingDisabled(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           INFO,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE, // Disabled
		piiMode:         SHOW_PII,       // Disabled
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	InfoWithFields("Test unmasked", map[string]any{
		"email":    "user@example.com",
		"password": "secret123",
	})

	output := buf.String()

	// Parse JSON to check fields
	var logEntry map[string]any
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		err := json.Unmarshal([]byte(lines[0]), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		fields, ok := logEntry["fields"].(map[string]any)
		if !ok {
			t.Fatalf("Fields not found in log entry")
		}

		// Check that data is NOT masked
		if fields["email"] != "user@example.com" {
			t.Errorf("Expected email to not be masked, got: %v", fields["email"])
		}

		if fields["password"] != "secret123" {
			t.Errorf("Expected password to not be masked, got: %v", fields["password"])
		}
	}
}

// TestPlainFormat tests plain text output format
func TestPlainFormat(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           INFO,
		writer:          &buf,
		format:          PLAIN_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
		component:       "test-app",
		version:         "v1.0.0",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	Info("Plain format test")
	output := buf.String()

	// Check plain format structure - account for ANSI color codes
	if !strings.Contains(output, "info") || !strings.Contains(output, " | ") {
		t.Errorf("Expected plain format info level not found: %s", output)
	}

	if !strings.Contains(output, "Plain format test") {
		t.Errorf("Expected message not found: %s", output)
	}
}

// TestNestedFieldMasking tests masking in nested objects
func TestNestedFieldMasking(t *testing.T) {
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

	// Test with nested data
	InfoWithFields("Nested test", map[string]any{
		"user": map[string]any{
			"email":    "user@example.com",
			"password": "secret123",
			"id":       456,
		},
		"safe_field": "safe_value",
	})

	output := buf.String()

	// Parse JSON to check nested fields
	var logEntry map[string]any
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		err := json.Unmarshal([]byte(lines[0]), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		fields, ok := logEntry["fields"].(map[string]any)
		if !ok {
			t.Fatalf("Fields not found in log entry")
		}

		userFields, ok := fields["user"].(map[string]any)
		if !ok {
			t.Fatalf("Nested user fields not found")
		}

		// Check nested masking
		if userFields["email"] != "***PII***" {
			t.Errorf("Expected nested email to be masked, got: %v", userFields["email"])
		}

		if userFields["password"] != "***MASKED***" {
			t.Errorf("Expected nested password to be masked, got: %v", userFields["password"])
		}

		// Check that safe fields are not masked
		if fields["safe_field"] != "safe_value" {
			t.Errorf("Expected safe field to not be masked, got: %v", fields["safe_field"])
		}
	}
}

// TestParseLogLevel tests log level parsing
func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", DEBUG},
		{"DEBUG", DEBUG},
		{"info", INFO},
		{"INFO", INFO},
		{"information", INFO},
		{"warn", WARN},
		{"WARN", WARN},
		{"warning", WARN},
		{"error", ERROR},
		{"ERROR", ERROR},
		{"invalid", INFO}, // Default fallback
		{"", INFO},        // Default fallback
	}

	for _, test := range tests {
		result := ParseLogLevel(test.input)
		if result != test.expected {
			t.Errorf("ParseLogLevel(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// TestLogLevelString tests log level string conversion
func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "debug"},
		{INFO, "info"},
		{WARN, "warn"},
		{ERROR, "error"},
		{LogLevel(999), "info"}, // Invalid level defaults to info
	}

	for _, test := range tests {
		result := test.level.String()
		if result != test.expected {
			t.Errorf("LogLevel(%v).String() = %q, expected %q", test.level, result, test.expected)
		}
	}
}

// TestConfigurationFunctions tests various configuration functions
func TestConfigurationFunctions(t *testing.T) {
	// Save original logger
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	// Reset to a clean state
	defaultLogger = &Logger{
		level:           INFO,
		writer:          os.Stdout,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	// Test SetComponent
	SetComponent("test-component")
	if defaultLogger.component != "test-component" {
		t.Errorf("SetComponent failed, got: %s", defaultLogger.component)
	}

	// Test SetVersion
	SetVersion("v2.0.0")
	if defaultLogger.version != "v2.0.0" {
		t.Errorf("SetVersion failed, got: %s", defaultLogger.version)
	}

	// Test SetLevel
	SetLevel("debug")
	if defaultLogger.level != DEBUG {
		t.Errorf("SetLevel failed, got: %v", defaultLogger.level)
	}

	// Test SetShowCaller
	SetShowCaller(true)
	if !defaultLogger.showCaller {
		t.Errorf("SetShowCaller failed, expected true")
	}

	// Test SetPlainFormat
	SetPlainFormat()
	if defaultLogger.format != PLAIN_FORMAT {
		t.Errorf("SetPlainFormat failed, got: %v", defaultLogger.format)
	}

	// Test SetJSONFormat
	SetJSONFormat()
	if defaultLogger.format != JSON_FORMAT {
		t.Errorf("SetJSONFormat failed, got: %v", defaultLogger.format)
	}

	// Test SetMaskString
	SetMaskString("[REDACTED]")
	if defaultLogger.maskString != "[REDACTED]" {
		t.Errorf("SetMaskString failed, got: %s", defaultLogger.maskString)
	}

	// Test SetPIIMaskString
	SetPIIMaskString("[PII_HIDDEN]")
	if defaultLogger.piiMaskString != "[PII_HIDDEN]" {
		t.Errorf("SetPIIMaskString failed, got: %s", defaultLogger.piiMaskString)
	}
}

// TestProductionAndDevelopmentModes tests the convenience mode functions
func TestProductionAndDevelopmentModes(t *testing.T) {
	// Save original logger
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	// Reset to a clean state
	defaultLogger = &Logger{
		level:           INFO,
		writer:          os.Stdout,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	// Test SetProductionMode
	SetProductionMode()
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Errorf("SetProductionMode should enable sensitive masking")
	}
	if defaultLogger.piiMode != MASK_PII {
		t.Errorf("SetProductionMode should enable PII masking")
	}
	if defaultLogger.format != JSON_FORMAT {
		t.Errorf("SetProductionMode should set JSON format")
	}
	if defaultLogger.level != INFO {
		t.Errorf("SetProductionMode should set INFO level")
	}
	if defaultLogger.showCaller != false {
		t.Errorf("SetProductionMode should disable caller info")
	}

	// Test SetDevelopmentMode
	SetDevelopmentMode()
	if defaultLogger.sensitiveMode != SHOW_SENSITIVE {
		t.Errorf("SetDevelopmentMode should disable sensitive masking")
	}
	if defaultLogger.piiMode != SHOW_PII {
		t.Errorf("SetDevelopmentMode should disable PII masking")
	}
	if defaultLogger.format != PLAIN_FORMAT {
		t.Errorf("SetDevelopmentMode should set PLAIN format")
	}
	if defaultLogger.level != DEBUG {
		t.Errorf("SetDevelopmentMode should set DEBUG level")
	}
	if defaultLogger.showCaller != true {
		t.Errorf("SetDevelopmentMode should enable caller info")
	}
}

// TestBackwardCompatibility tests the JSON and Plain functions
func TestBackwardCompatibility(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           INFO,
		writer:          &buf,
		format:          PLAIN_FORMAT, // Set to plain, but JSON/Plain functions should override
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Test JSON function (should force JSON output)
	JSON("info", "JSON test message")

	output := buf.String()
	if !strings.Contains(output, `"msg":"JSON test message"`) {
		t.Errorf("JSON function should produce JSON output: %s", output)
	}

	// Clear buffer
	buf.Reset()

	// Test Plain function (should force plain output)
	Plain("info", "Plain test message")

	output = buf.String()
	if !strings.Contains(output, "info") || !strings.Contains(output, " | ") {
		t.Errorf("Plain function should produce plain output: %s", output)
	}
	if !strings.Contains(output, "Plain test message") {
		t.Errorf("Plain function should contain message: %s", output)
	}
}
