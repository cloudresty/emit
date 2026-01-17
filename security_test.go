package emit

import (
	"bytes"
	"strings"
	"testing"
)

// TestSensitiveFieldMasking tests that sensitive fields are properly masked
func TestSensitiveFieldMasking(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Test with a map containing sensitive fields
	Info.Field("Login attempt",
		NewFields().
			String("password", "secret123").
			String("api_key", "sk-1234567890").
			String("username", "testuser"))

	output := buf.String()

	// Password and API key should be masked
	if strings.Contains(output, "secret123") {
		t.Error("Password should be masked")
	}
	if strings.Contains(output, "sk-1234567890") {
		t.Error("API key should be masked")
	}

	// Username should be visible
	if !strings.Contains(output, "testuser") {
		t.Error("Username should be visible")
	}
}

// TestPIIFieldMasking tests that PII fields are properly masked
func TestPIIFieldMasking(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	Info.Field("User registration",
		NewFields().
			String("email", "user@example.com").
			String("phone", "555-123-4567").
			String("status", "active"))

	output := buf.String()

	// Email and phone should be masked as PII
	if strings.Contains(output, "user@example.com") {
		t.Error("Email should be masked as PII")
	}
	if strings.Contains(output, "555-123-4567") {
		t.Error("Phone should be masked as PII")
	}

	// Status should be visible
	if !strings.Contains(output, "active") {
		t.Error("Status should be visible")
	}
}

// TestBothMaskingModes tests sensitive and PII masking together
func TestBothMaskingModes(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
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

	Info.Field("Combined test",
		NewFields().
			String("password", "secret").
			String("email", "test@test.com").
			String("message", "Hello"))

	output := buf.String()

	// Both should be masked
	if strings.Contains(output, "secret") {
		t.Error("Password should be masked")
	}
	if strings.Contains(output, "test@test.com") {
		t.Error("Email should be masked")
	}

	// Message should be visible
	if !strings.Contains(output, "Hello") {
		t.Error("Message should be visible")
	}
}

// TestNoMasking tests that all data is visible when masking is disabled
func TestNoMasking(t *testing.T) {
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

	Info.Field("No masking test",
		NewFields().
			String("password", "visible_password").
			String("email", "visible@email.com"))

	output := buf.String()

	// Both should be visible
	if !strings.Contains(output, "visible_password") {
		t.Error("Password should be visible when masking is off")
	}
	if !strings.Contains(output, "visible@email.com") {
		t.Error("Email should be visible when masking is off")
	}
}

