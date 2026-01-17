package emit

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestJSONFormatter tests JSON format output
func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
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

	Info.Msg("JSON test message")

	output := buf.String()

	// Verify JSON structure
	if !strings.Contains(output, `"timestamp":`) {
		t.Error("JSON output should contain timestamp")
	}
	if !strings.Contains(output, `"level":"info"`) {
		t.Error("JSON output should contain level")
	}
	if !strings.Contains(output, `"message":"JSON test message"`) {
		t.Error("JSON output should contain message")
	}
	if !strings.Contains(output, `"component":"test-component"`) {
		t.Error("JSON output should contain component")
	}
	if !strings.Contains(output, `"version":"1.0.0"`) {
		t.Error("JSON output should contain version")
	}
}

// TestPlainFormatter tests plain text format output
func TestPlainFormatter(t *testing.T) {
	var buf bytes.Buffer

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

	Info.Msg("Plain test message")

	output := buf.String()

	// Plain format should contain level (lowercase) and message
	if !strings.Contains(output, "info") {
		t.Error("Plain output should contain info level")
	}
	if !strings.Contains(output, "Plain test message") {
		t.Error("Plain output should contain message")
	}
}

// TestDynamicBufferPath tests the dynamic buffer allocation path
// This is triggered when message + fields exceed the static buffer size
func TestDynamicBufferPath(t *testing.T) {
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

	// Create a very long message to trigger dynamic buffer path
	longMessage := strings.Repeat("This is a very long message. ", 50)

	// Add many fields to ensure we exceed the buffer
	Info.StructuredFields(longMessage,
		ZString("field1", strings.Repeat("value", 100)),
		ZString("field2", strings.Repeat("data", 100)),
		ZString("field3", strings.Repeat("content", 100)),
		ZString("field4", strings.Repeat("text", 100)),
		ZString("field5", strings.Repeat("info", 100)),
		ZInt("count", 42),
		ZInt64("big_num", 9223372036854775807),
		ZFloat64("rate", 3.14159),
		ZBool("flag", true),
	)

	output := buf.String()

	// Verify the output is valid JSON and contains expected content
	if !strings.Contains(output, longMessage[:50]) {
		t.Error("Output should contain the long message")
	}
	if !strings.Contains(output, `"field1":`) {
		t.Error("Output should contain field1")
	}
	if !strings.Contains(output, `"count":42`) {
		t.Error("Output should contain count field")
	}
}

// TestStructuredFieldsWithManyFields tests with >4 fields to trigger estimation path
func TestStructuredFieldsWithManyFields(t *testing.T) {
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

	testTime := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)

	// Use more than 4 fields to trigger the size estimation branch
	Info.StructuredFields("Many fields test",
		ZString("str1", "value1"),
		ZString("str2", "value2"),
		ZInt("int1", 100),
		ZInt("int2", 200),
		ZInt64("int64_1", 1234567890123),
		ZFloat64("float1", 1.23),
		ZBool("bool1", true),
		ZTime("time1", testTime),
		ZDuration("dur1", 5*time.Minute),
	)

	output := buf.String()

	expectedFields := []string{
		`"str1":"value1"`,
		`"str2":"value2"`,
		`"int1":100`,
		`"int2":200`,
		`"int64_1":1234567890123`,
		`"float1":1.23`,
		`"bool1":true`,
		`"time1":"2024-06-15T14:30:00Z"`,
		`"dur1":"5m0s"`,
	}

	for _, expected := range expectedFields {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected %s not found in output: %s", expected, output)
		}
	}
}

// TestJSONEscaping tests proper JSON string escaping
func TestJSONEscaping(t *testing.T) {
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

	// Test with special characters that need escaping
	Info.StructuredFields("Escape test",
		ZString("quotes", `He said "hello"`),
		ZString("backslash", `path\to\file`),
		ZString("newline", "line1\nline2"),
		ZString("tab", "col1\tcol2"),
	)

	output := buf.String()

	// Check that special chars are escaped
	if !strings.Contains(output, `\"hello\"`) {
		t.Error("Quotes should be escaped")
	}
	if !strings.Contains(output, `\\`) {
		t.Error("Backslashes should be escaped")
	}
}

// TestEscapeJSONString tests the escapeJSONString helper
func TestEscapeJSONString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with space", "with space"},
		{`quote"here`, `quote\"here`},
		{"new\nline", `new\nline`},
	}

	for _, tt := range tests {
		dst := make([]byte, 100)
		n := escapeJSONString(dst, tt.input)
		result := string(dst[:n])
		if result != tt.expected {
			t.Errorf("escapeJSONString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestPlainFormatterWithFields tests plain format with structured fields
func TestPlainFormatterWithFields(t *testing.T) {
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

	Info.Field("Plain with fields",
		NewFields().
			String("key", "value").
			Int("count", 42))

	output := buf.String()

	if !strings.Contains(output, "Plain with fields") {
		t.Error("Plain output should contain message")
	}
}

// TestAllLogLevelsJSON tests all log levels produce correct JSON
func TestAllLogLevelsJSON(t *testing.T) {
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

	Debug.Msg("Debug test")
	Info.Msg("Info test")
	Warn.Msg("Warn test")
	Error.Msg("Error test")

	output := buf.String()

	expectedLevels := []string{"debug", "info", "warn", "error"}
	for _, level := range expectedLevels {
		expected := `"level":"` + level + `"`
		if !strings.Contains(output, expected) {
			t.Errorf("Expected level %s in output, got: %s", level, output)
		}
	}
}

// TestPlainFormatWithFields tests plain format with fields
func TestPlainFormatWithFields(t *testing.T) {
	var buf bytes.Buffer

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

	Info.Field("Plain with fields", NewFields().String("key", "value").Int("count", 42))

	output := buf.String()
	if !strings.Contains(output, "Plain with fields") {
		t.Error("Plain format should log message with fields")
	}
}

// TestPlainFormatLongMessage tests plain format with long messages
func TestPlainFormatLongMessage(t *testing.T) {
	var buf bytes.Buffer

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

	longMessage := strings.Repeat("Long plain message. ", 100)
	Info.Msg(longMessage)

	output := buf.String()
	if !strings.Contains(output, "Long plain message") {
		t.Error("Plain format should handle long messages")
	}
}

// TestJSONFormatWithCaller tests JSON format with caller info enabled
func TestJSONFormatWithCaller(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		component:       "test-component",
		version:         "1.0.0",
		showCaller:      true,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	Info.Field("Caller test", NewFields().String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, "Caller test") {
		t.Error("JSON format with caller should log message")
	}
	// Caller info should be present
	if !strings.Contains(output, "file") {
		t.Error("JSON format with caller should include file info")
	}
}

// TestJSONFormatWithComponentAndVersion tests JSON format with component and version
func TestJSONFormatWithComponentAndVersion(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		component:       "my-service",
		version:         "2.0.0",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	Info.Field("Component test", NewFields().String("key", "value"))

	output := buf.String()
	if !strings.Contains(output, "my-service") {
		t.Error("JSON format should include component")
	}
	if !strings.Contains(output, "2.0.0") {
		t.Error("JSON format should include version")
	}
}

// TestDynamicBufferWithAllFieldTypes tests dynamic buffer with all field types
func TestDynamicBufferWithAllFieldTypes(t *testing.T) {
	var buf bytes.Buffer

	testLogger := &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
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

	// Create a very long message with all field types to trigger dynamic buffer
	longMessage := strings.Repeat("Very long message content for dynamic buffer testing. ", 300)

	Info.StructuredFields(longMessage,
		ZString("str_field", strings.Repeat("long_value", 200)),
		ZInt("int_field", 42),
		ZInt64("int64_field", 9223372036854775807),
		ZFloat64("float_field", 3.14159),
		ZBool("bool_field", true),
		ZTime("time_field", time.Now()),
		ZDuration("duration_field", 5*time.Second),
	)

	output := buf.String()
	if !strings.Contains(output, "Very long message") {
		t.Error("Dynamic buffer should handle very long messages with all field types")
	}
	if !strings.Contains(output, `"int_field":42`) {
		t.Error("Dynamic buffer should include int field")
	}
	if !strings.Contains(output, `"bool_field":true`) {
		t.Error("Dynamic buffer should include bool field")
	}
}

// TestPlainFormatAllLevels tests plain format with all log levels
func TestPlainFormatAllLevels(t *testing.T) {
	var buf bytes.Buffer

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

	Debug.Msg("Debug plain message")
	Info.Msg("Info plain message")
	Warn.Msg("Warn plain message")
	Error.Msg("Error plain message")

	output := buf.String()
	if !strings.Contains(output, "debug") {
		t.Error("Plain format should include debug level")
	}
	if !strings.Contains(output, "info") {
		t.Error("Plain format should include info level")
	}
	if !strings.Contains(output, "warn") {
		t.Error("Plain format should include warn level")
	}
	if !strings.Contains(output, "error") {
		t.Error("Plain format should include error level")
	}
}

// TestLongZStringNoPanic tests that very long ZString values don't cause panic
func TestLongZStringNoPanic(t *testing.T) {
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

	// This previously could cause a panic - short message (<200 chars),
	// only 1 field (<4 fields), but very long value (>1024 chars)
	longValue := strings.Repeat("x", 2000)

	// Should NOT panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Long ZString value caused panic: %v", r)
		}
	}()

	Info.StructuredFields("Short message", ZString("data", longValue))

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty string")
	}
}

// TestVeryLongZStringWithEscaping tests very long strings that need JSON escaping
func TestVeryLongZStringWithEscaping(t *testing.T) {
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

	// Long string with characters that need escaping (quotes, backslashes, newlines)
	longValue := strings.Repeat("hello\"world\\test\n", 200)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Long ZString with escaping caused panic: %v", r)
		}
	}()

	Info.StructuredFields("Escape test", ZString("data", longValue))

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty string")
	}
	// Verify escaping happened
	if !strings.Contains(output, `\\`) {
		t.Error("Expected escaped backslash in output")
	}
}

// TestMultipleLongZStrings tests multiple long ZString fields
func TestMultipleLongZStrings(t *testing.T) {
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

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Multiple long ZStrings caused panic: %v", r)
		}
	}()

	Info.StructuredFields("Multi long fields",
		ZString("field1", strings.Repeat("a", 500)),
		ZString("field2", strings.Repeat("b", 500)),
		ZString("field3", strings.Repeat("c", 500)),
	)

	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected output, got empty string")
	}
}
