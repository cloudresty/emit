package emit

import (
	"bytes"
	"strings"
	"testing"
)

// TestLegacyEmitFunctions tests the legacy emit functions
func TestLegacyEmitFunctions(t *testing.T) {
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

	// Test InfoMsg
	InfoMsg("Info message")
	if !strings.Contains(buf.String(), "Info message") {
		t.Error("InfoMsg should log info message")
	}

	buf.Reset()

	// Test ErrorMsg
	ErrorMsg("Error message")
	if !strings.Contains(buf.String(), "Error message") {
		t.Error("ErrorMsg should log error message")
	}

	buf.Reset()

	// Test WarnMsg
	WarnMsg("Warn message")
	if !strings.Contains(buf.String(), "Warn message") {
		t.Error("WarnMsg should log warn message")
	}

	buf.Reset()

	// Test DebugMsg
	DebugMsg("Debug message")
	if !strings.Contains(buf.String(), "Debug message") {
		t.Error("DebugMsg should log debug message")
	}
}

// TestInfoWithFields tests InfoWithFields legacy function
func TestInfoWithFields(t *testing.T) {
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

	InfoWithFields("Test with fields", map[string]any{
		"key": "value",
	})

	output := buf.String()
	if !strings.Contains(output, "Test with fields") {
		t.Error("InfoWithFields should log message")
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Error("InfoWithFields should include fields")
	}
}

// TestLogFunction tests the Log function
func TestLogFunction(t *testing.T) {
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

	// Log takes optional params as component and version, not key-value
	Log("info", "Log function test", "my-component", "1.0.0")

	output := buf.String()
	if !strings.Contains(output, "Log function test") {
		t.Error("Log should log message")
	}
	if !strings.Contains(output, `"level":"info"`) {
		t.Error("Log should include level")
	}
}

// TestJSONFormatFunction tests the JSON format function
func TestJSONFormatFunction(t *testing.T) {
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

	JSON("info", "JSON test", "field", "data")

	output := buf.String()
	if !strings.Contains(output, "JSON test") {
		t.Error("JSON should log message")
	}
	if !strings.Contains(output, `"level":"info"`) {
		t.Error("JSON should produce JSON output")
	}
}

// TestPlainFormatFunction tests the Plain format function
func TestPlainFormatFunction(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          PLAIN_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	Plain("info", "Plain test", "field", "data")

	output := buf.String()
	if !strings.Contains(output, "Plain test") {
		t.Error("Plain should log message")
	}
}

