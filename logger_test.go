package emit

import (
	"bytes"
	"strings"
	"testing"
)

// TestLoggerStructuredMethods tests the Logger's structured logging methods
func TestLoggerStructuredMethods(t *testing.T) {
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

	// Test package-level InfoStructured
	InfoStructured("Info structured test", ZString("key", "value"))
	if !strings.Contains(buf.String(), "Info structured test") {
		t.Error("InfoStructured should log message")
	}

	buf.Reset()

	// Test package-level DebugStructured
	DebugStructured("Debug structured test", ZInt("count", 42))
	if !strings.Contains(buf.String(), "Debug structured test") {
		t.Error("DebugStructured should log message")
	}

	buf.Reset()

	// Test package-level WarnStructured
	WarnStructured("Warn structured test", ZBool("flag", true))
	if !strings.Contains(buf.String(), "Warn structured test") {
		t.Error("WarnStructured should log message")
	}

	buf.Reset()

	// Test package-level ErrorStructured
	ErrorStructured("Error structured test", ZFloat64("rate", 3.14))
	if !strings.Contains(buf.String(), "Error structured test") {
		t.Error("ErrorStructured should log message")
	}
}

// TestLoggerInstanceMethods tests the Logger instance methods
func TestLoggerInstanceMethods(t *testing.T) {
	var buf bytes.Buffer

	logger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	// Test instance InfoStructured
	logger.InfoStructured("Logger info test", ZString("instance", "true"))
	if !strings.Contains(buf.String(), "Logger info test") {
		t.Error("Logger.InfoStructured should log message")
	}

	buf.Reset()

	// Test instance DebugStructured
	logger.DebugStructured("Logger debug test", ZInt("num", 1))
	if !strings.Contains(buf.String(), "Logger debug test") {
		t.Error("Logger.DebugStructured should log message")
	}

	buf.Reset()

	// Test instance WarnStructured
	logger.WarnStructured("Logger warn test", ZBool("warning", true))
	if !strings.Contains(buf.String(), "Logger warn test") {
		t.Error("Logger.WarnStructured should log message")
	}

	buf.Reset()

	// Test instance ErrorStructured
	logger.ErrorStructured("Logger error test", ZString("err", "test"))
	if !strings.Contains(buf.String(), "Logger error test") {
		t.Error("Logger.ErrorStructured should log message")
	}
}

// TestLoggerPlainSizeEstimation tests plain format size estimation
func TestLoggerPlainSizeEstimation(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger with component and version to test size estimation
	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          PLAIN_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		component:       "test-component",
		version:         "1.0.0",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Log a long message to trigger size estimation
	longMessage := strings.Repeat("This is a very long message. ", 20)
	Info.Msg(longMessage)

	output := buf.String()
	if !strings.Contains(output, "This is a very long message") {
		t.Error("Plain format should log long messages correctly")
	}
}

