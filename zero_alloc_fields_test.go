package emit

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestZFieldConstructors tests all ZField constructor functions
func TestZFieldConstructors(t *testing.T) {
	// Test ZString
	strField := ZString("key", "value")
	if strField.Key != "key" || strField.Value != "value" {
		t.Errorf("ZString failed: got key=%s, value=%s", strField.Key, strField.Value)
	}

	// Test ZInt
	intField := ZInt("count", 42)
	if intField.Key != "count" || intField.Value != 42 {
		t.Errorf("ZInt failed: got key=%s, value=%d", intField.Key, intField.Value)
	}

	// Test ZInt64
	int64Field := ZInt64("big_number", 9223372036854775807)
	if int64Field.Key != "big_number" || int64Field.Value != 9223372036854775807 {
		t.Errorf("ZInt64 failed: got key=%s, value=%d", int64Field.Key, int64Field.Value)
	}

	// Test ZFloat64
	floatField := ZFloat64("rate", 3.14159)
	if floatField.Key != "rate" || floatField.Value != 3.14159 {
		t.Errorf("ZFloat64 failed: got key=%s, value=%f", floatField.Key, floatField.Value)
	}

	// Test ZBool
	boolField := ZBool("enabled", true)
	if boolField.Key != "enabled" || boolField.Value != true {
		t.Errorf("ZBool failed: got key=%s, value=%t", boolField.Key, boolField.Value)
	}

	// Test ZTime
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timeField := ZTime("timestamp", testTime)
	if timeField.Key != "timestamp" || !timeField.Value.Equal(testTime) {
		t.Errorf("ZTime failed: got key=%s, value=%v", timeField.Key, timeField.Value)
	}

	// Test ZDuration
	durField := ZDuration("elapsed", 5*time.Second)
	if durField.Key != "elapsed" || durField.Value != 5*time.Second {
		t.Errorf("ZDuration failed: got key=%s, value=%v", durField.Key, durField.Value)
	}
}

// TestZFieldsInStructuredLogging tests that all ZField types work in structured logging
func TestZFieldsInStructuredLogging(t *testing.T) {
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

	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	// Test all field types together
	Info.StructuredFields("Complete field test",
		ZString("string_field", "test_value"),
		ZInt("int_field", 42),
		ZInt64("int64_field", 9223372036854775807),
		ZFloat64("float_field", 3.14159),
		ZBool("bool_true", true),
		ZBool("bool_false", false),
		ZTime("time_field", testTime),
		ZDuration("duration_field", 2*time.Hour+30*time.Minute),
	)

	output := buf.String()

	// Verify all fields are present
	expectedFields := []string{
		`"string_field":"test_value"`,
		`"int_field":42`,
		`"int64_field":9223372036854775807`,
		`"float_field":3.14159`,
		`"bool_true":true`,
		`"bool_false":false`,
		`"time_field":"2024-01-15T10:30:00Z"`,
		`"duration_field":"2h30m0s"`,
	}

	for _, expected := range expectedFields {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected field %s not found in output: %s", expected, output)
		}
	}
}

// TestZFieldSensitivity tests IsSensitive and IsPII methods
func TestZFieldSensitivity(t *testing.T) {
	// Test sensitive string fields
	passwordField := ZString("password", "secret123")
	if !passwordField.IsSensitive() {
		t.Error("password field should be sensitive")
	}

	// Test PII string fields
	emailField := ZString("email", "user@example.com")
	if !emailField.IsPII() {
		t.Error("email field should be PII")
	}

	// Test non-sensitive fields
	normalField := ZString("status", "active")
	if normalField.IsSensitive() {
		t.Error("status field should not be sensitive")
	}
	if normalField.IsPII() {
		t.Error("status field should not be PII")
	}

	// Test that numeric fields are never sensitive/PII
	intField := ZInt("count", 42)
	if intField.IsSensitive() || intField.IsPII() {
		t.Error("IntZField should never be sensitive or PII")
	}

	int64Field := ZInt64("id", 123456789)
	if int64Field.IsSensitive() || int64Field.IsPII() {
		t.Error("Int64ZField should never be sensitive or PII")
	}

	floatField := ZFloat64("rate", 3.14)
	if floatField.IsSensitive() || floatField.IsPII() {
		t.Error("Float64ZField should never be sensitive or PII")
	}

	boolField := ZBool("active", true)
	if boolField.IsSensitive() || boolField.IsPII() {
		t.Error("BoolZField should never be sensitive or PII")
	}

	timeField := ZTime("created", time.Now())
	if timeField.IsSensitive() || timeField.IsPII() {
		t.Error("TimeZField should never be sensitive or PII")
	}

	durField := ZDuration("elapsed", time.Second)
	if durField.IsSensitive() || durField.IsPII() {
		t.Error("DurationZField should never be sensitive or PII")
	}
}

// TestZFieldMaskingInOutput tests that sensitive fields are properly masked
func TestZFieldMaskingInOutput(t *testing.T) {
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

	Info.StructuredFields("Sensitive data test",
		ZString("password", "secret123"),
		ZString("email", "user@example.com"),
		ZString("normal", "visible"),
	)

	output := buf.String()

	// Password should be masked
	if strings.Contains(output, "secret123") {
		t.Error("Password should be masked but was visible in output")
	}
	if !strings.Contains(output, "***MASKED***") {
		t.Error("Masked placeholder not found in output")
	}

	// Normal field should be visible
	if !strings.Contains(output, "visible") {
		t.Error("Normal field should be visible")
	}
}

// TestAllLogLevelsWithStructuredFields tests structured fields with all log levels
func TestAllLogLevelsWithStructuredFields(t *testing.T) {
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

	// Test Debug.StructuredFields
	Debug.StructuredFields("Debug message",
		ZString("level", "debug"),
		ZInt("code", 1))

	// Test Info.StructuredFields (already tested elsewhere, but for completeness)
	Info.StructuredFields("Info message",
		ZString("level", "info"),
		ZInt("code", 2))

	// Test Warn.StructuredFields
	Warn.StructuredFields("Warn message",
		ZString("level", "warn"),
		ZInt("code", 3))

	// Test Error.StructuredFields
	Error.StructuredFields("Error message",
		ZString("level", "error"),
		ZInt("code", 4))

	output := buf.String()

	expectedLevels := []string{
		`"level":"debug"`,
		`"level":"info"`,
		`"level":"warn"`,
		`"level":"error"`,
	}

	for _, expected := range expectedLevels {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected %s not found in output", expected)
		}
	}
}

// TestZFieldWriteToEncoder tests the WriteToEncoder methods
func TestZFieldWriteToEncoder(t *testing.T) {
	enc := &ZeroAllocEncoder{}
	enc.buf = make([]byte, 0, 1024)

	// Use a key that's not a PII field to avoid masking
	strField := ZString("status", "active")
	strField.WriteToEncoder(enc)

	// Test IntZField.WriteToEncoder
	intField := ZInt("count", 42)
	intField.WriteToEncoder(enc)

	// Test Int64ZField.WriteToEncoder
	int64Field := ZInt64("big", 9223372036854775807)
	int64Field.WriteToEncoder(enc)

	// Test Float64ZField.WriteToEncoder
	floatField := ZFloat64("rate", 3.14)
	floatField.WriteToEncoder(enc)

	// Test BoolZField.WriteToEncoder
	boolField := ZBool("flag", true)
	boolField.WriteToEncoder(enc)

	// Test TimeZField.WriteToEncoder
	testTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timeField := ZTime("created", testTime)
	timeField.WriteToEncoder(enc)

	// Test DurationZField.WriteToEncoder - encoder stores as nanoseconds
	durField := ZDuration("elapsed", 5*time.Second)
	durField.WriteToEncoder(enc)

	output := string(enc.buf)

	// The encoder stores duration as nanoseconds (5 seconds = 5000000000 ns)
	expectedParts := []string{
		`"status":"active"`,
		`"count":42`,
		`"big":9223372036854775807`,
		`"rate":3.14`,
		`"flag":true`,
		`"created":"2024-01-15T10:30:00Z"`,
		`"elapsed":5000000000`, // Duration stored as nanoseconds
	}

	for _, expected := range expectedParts {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected %s not found in encoder output: %s", expected, output)
		}
	}
}

// TestStringZFieldMasking tests StringZField masking behavior
func TestStringZFieldMasking(t *testing.T) {
	// Test with a sensitive field
	sensitiveField := ZString("password", "secret123")
	if !sensitiveField.IsSensitive() {
		t.Error("password should be sensitive")
	}

	// Test with a PII field
	piiField := ZString("email", "test@example.com")
	if !piiField.IsPII() {
		t.Error("email should be PII")
	}

	// Test with a non-sensitive, non-PII field
	normalField := ZString("status", "active")
	if normalField.IsSensitive() {
		t.Error("status should not be sensitive")
	}
	if normalField.IsPII() {
		t.Error("status should not be PII")
	}
}

// TestEncoderWriteStringEscaping tests string escaping in encoder
func TestEncoderWriteStringEscaping(t *testing.T) {
	enc := &ZeroAllocEncoder{}
	enc.buf = make([]byte, 0, 1024)

	// Test string with special characters that need escaping
	specialField := ZString("message", "Hello\nWorld\t\"quoted\"\\backslash")
	specialField.WriteToEncoder(enc)

	output := string(enc.buf)

	// Should contain escaped characters
	if !strings.Contains(output, `\n`) {
		t.Error("Newline should be escaped")
	}
	if !strings.Contains(output, `\t`) {
		t.Error("Tab should be escaped")
	}
	if !strings.Contains(output, `\"`) {
		t.Error("Quote should be escaped")
	}
	if !strings.Contains(output, `\\`) {
		t.Error("Backslash should be escaped")
	}
}

// TestEncoderWriteStringControlChars tests control character escaping
func TestEncoderWriteStringControlChars(t *testing.T) {
	enc := &ZeroAllocEncoder{}
	enc.buf = make([]byte, 0, 1024)

	// Test string with control characters (ASCII < 32)
	controlField := ZString("data", "test\x01\x02\x03end")
	controlField.WriteToEncoder(enc)

	output := string(enc.buf)

	// Should contain unicode escape sequences for control chars
	if !strings.Contains(output, `\u00`) {
		t.Error("Control characters should be escaped as unicode")
	}
}

// TestEncoderWriteStringUnicode tests unicode character handling
func TestEncoderWriteStringUnicode(t *testing.T) {
	enc := &ZeroAllocEncoder{}
	enc.buf = make([]byte, 0, 1024)

	// Test string with unicode characters
	unicodeField := ZString("greeting", "Hello 世界 🌍")
	unicodeField.WriteToEncoder(enc)

	output := string(enc.buf)

	// Should contain the unicode characters
	if !strings.Contains(output, "世界") {
		t.Error("Unicode characters should be preserved")
	}
}

// TestStringZFieldMaskingInEncoder tests that sensitive/PII fields are masked in encoder
func TestStringZFieldMaskingInEncoder(t *testing.T) {
	// Test sensitive field masking
	enc := &ZeroAllocEncoder{}
	enc.buf = make([]byte, 0, 1024)

	sensitiveField := ZString("password", "secret123")
	sensitiveField.WriteToEncoder(enc)

	output := string(enc.buf)
	if strings.Contains(output, "secret123") {
		t.Error("Sensitive field value should be masked")
	}
	if !strings.Contains(output, "***MASKED***") {
		t.Error("Sensitive field should show ***MASKED***")
	}

	// Test PII field masking
	enc2 := &ZeroAllocEncoder{}
	enc2.buf = make([]byte, 0, 1024)

	piiField := ZString("email", "test@example.com")
	piiField.WriteToEncoder(enc2)

	output2 := string(enc2.buf)
	if strings.Contains(output2, "test@example.com") {
		t.Error("PII field value should be masked")
	}
	if !strings.Contains(output2, "***PII***") {
		t.Error("PII field should show ***PII***")
	}
}

// TestAllSensitiveFieldKeys tests all sensitive field key patterns
func TestAllSensitiveFieldKeys(t *testing.T) {
	sensitiveKeys := []string{
		"password", "secret", "token", "api_key", "private_key",
		"auth", "credential", "session", "jwt", "bearer",
	}

	for _, key := range sensitiveKeys {
		field := ZString(key, "test_value")
		if !field.IsSensitive() {
			t.Errorf("Field with key '%s' should be sensitive", key)
		}
	}
}

// TestAllPIIFieldKeys tests all PII field key patterns
func TestAllPIIFieldKeys(t *testing.T) {
	piiKeys := []string{
		"email", "phone", "name", "address", "ssn",
		"user_email", "full_name", "credit_card", "passport",
	}

	for _, key := range piiKeys {
		field := ZString(key, "test_value")
		if !field.IsPII() {
			t.Errorf("Field with key '%s' should be PII", key)
		}
	}
}

