package emit

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestSetComponent tests component configuration
func TestSetComponent(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	var buf bytes.Buffer
	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	SetComponent("my-service")
	Info.Msg("Test message")

	output := buf.String()
	if !strings.Contains(output, `"component":"my-service"`) {
		t.Errorf("Expected component in output, got: %s", output)
	}
}

// TestSetVersion tests version configuration
func TestSetVersion(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	var buf bytes.Buffer
	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	SetVersion("2.0.0")
	Info.Msg("Test message")

	output := buf.String()
	if !strings.Contains(output, `"version":"2.0.0"`) {
		t.Errorf("Expected version in output, got: %s", output)
	}
}

// TestSetLevelString tests setting log level via string
func TestSetLevelString(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	var buf bytes.Buffer
	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	// Test setting to error level - debug and info should be filtered
	SetLevel("error")
	Debug.Msg("Debug message")
	Info.Msg("Info message")
	Error.Msg("Error message")

	output := buf.String()
	if strings.Contains(output, "Debug message") {
		t.Error("Debug message should be filtered out")
	}
	if strings.Contains(output, "Info message") {
		t.Error("Info message should be filtered out")
	}
	if !strings.Contains(output, "Error message") {
		t.Error("Error message should be present")
	}
}

// TestSetFormat tests format switching
func TestSetFormat(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	var buf bytes.Buffer
	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	// Test JSON format
	SetFormat("json")
	Info.Msg("JSON message")
	jsonOutput := buf.String()
	if !strings.Contains(jsonOutput, `"level":"info"`) {
		t.Error("JSON format should produce JSON output")
	}

	// Test Plain format
	buf.Reset()
	SetFormat("plain")
	Info.Msg("Plain message")
	plainOutput := buf.String()
	if strings.Contains(plainOutput, `"level"`) {
		t.Error("Plain format should not produce JSON")
	}

	// Test invalid format (should be ignored)
	SetFormat("invalid")
}

// TestSetPlainAndJSONFormat tests format helper methods
func TestSetPlainAndJSONFormat(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	var buf bytes.Buffer
	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	SetPlainFormat()
	if defaultLogger.format != PLAIN_FORMAT {
		t.Error("SetPlainFormat should set PLAIN_FORMAT")
	}

	SetJSONFormat()
	if defaultLogger.format != JSON_FORMAT {
		t.Error("SetJSONFormat should set JSON_FORMAT")
	}
}

// TestSetShowCaller tests caller information display
func TestSetShowCaller(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	var buf bytes.Buffer
	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          &buf,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	SetShowCaller(true)
	if !defaultLogger.showCaller {
		t.Error("SetShowCaller(true) should enable caller display")
	}

	SetShowCaller(false)
	if defaultLogger.showCaller {
		t.Error("SetShowCaller(false) should disable caller display")
	}
}

// TestSetOutput tests custom output writer
func TestSetOutput(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	var buf bytes.Buffer
	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          os.Stdout,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	SetOutput(&buf)
	Info.Msg("Test to custom writer")

	if buf.Len() == 0 {
		t.Error("SetOutput should redirect logs to custom writer")
	}
}

// TestSetOutputToDiscard tests discarding output
func TestSetOutputToDiscardFunc(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          os.Stdout,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
	}

	SetOutputToDiscard()
	if defaultLogger.writer != io.Discard {
		t.Error("SetOutputToDiscard should set writer to io.Discard")
	}
}

// TestSensitiveConfiguration tests sensitive data masking configuration
func TestSensitiveConfiguration(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          io.Discard,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: []string{},
		piiFields:       []string{},
	}

	// Test SetSensitiveMode with string
	SetSensitiveMode("mask")
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Error("SetSensitiveMode should update sensitive mode")
	}

	// Test ShowSensitiveData
	ShowSensitiveData()
	if defaultLogger.sensitiveMode != SHOW_SENSITIVE {
		t.Error("ShowSensitiveData should set SHOW_SENSITIVE mode")
	}

	// Test MaskSensitiveData
	MaskSensitiveData()
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Error("MaskSensitiveData should set MASK_SENSITIVE mode")
	}

	// Test SetMaskString
	SetMaskString("***HIDDEN***")
	if defaultLogger.maskString != "***HIDDEN***" {
		t.Error("SetMaskString should update mask string")
	}

	// Test AddSensitiveField
	AddSensitiveField("api_key")
	found := false
	for _, f := range defaultLogger.sensitiveFields {
		if f == "api_key" {
			found = true
			break
		}
	}
	if !found {
		t.Error("AddSensitiveField should add field to sensitive list")
	}

	// Test SetSensitiveFields
	SetSensitiveFields([]string{"secret", "token"})
	hasSecret := false
	hasToken := false
	for _, f := range defaultLogger.sensitiveFields {
		if f == "secret" {
			hasSecret = true
		}
		if f == "token" {
			hasToken = true
		}
	}
	if !hasSecret || !hasToken {
		t.Error("SetSensitiveFields should set sensitive fields list")
	}
}

// TestPIIConfiguration tests PII masking configuration
func TestPIIConfiguration(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          io.Discard,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: []string{},
		piiFields:       []string{},
	}

	// Test SetPIIMode with string
	SetPIIMode("mask")
	if defaultLogger.piiMode != MASK_PII {
		t.Error("SetPIIMode should update PII mode")
	}

	// Test ShowPIIData
	ShowPIIData()
	if defaultLogger.piiMode != SHOW_PII {
		t.Error("ShowPIIData should set SHOW_PII mode")
	}

	// Test MaskPIIData
	MaskPIIData()
	if defaultLogger.piiMode != MASK_PII {
		t.Error("MaskPIIData should set MASK_PII mode")
	}

	// Test SetPIIMaskString
	SetPIIMaskString("***REDACTED***")
	if defaultLogger.piiMaskString != "***REDACTED***" {
		t.Error("SetPIIMaskString should update PII mask string")
	}

	// Test AddPIIField
	AddPIIField("ssn")
	found := false
	for _, f := range defaultLogger.piiFields {
		if f == "ssn" {
			found = true
			break
		}
	}
	if !found {
		t.Error("AddPIIField should add field to PII list")
	}

	// Test SetPIIFields
	SetPIIFields([]string{"phone", "address"})
	hasPhone := false
	hasAddress := false
	for _, f := range defaultLogger.piiFields {
		if f == "phone" {
			hasPhone = true
		}
		if f == "address" {
			hasAddress = true
		}
	}
	if !hasPhone || !hasAddress {
		t.Error("SetPIIFields should set PII fields list")
	}
}

// TestSetAllMasking tests combined masking toggle
func TestSetAllMasking(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          io.Discard,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: []string{},
		piiFields:       []string{},
	}

	// Test enabling all masking
	SetAllMasking(true)
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Error("SetAllMasking(true) should enable sensitive masking")
	}
	if defaultLogger.piiMode != MASK_PII {
		t.Error("SetAllMasking(true) should enable PII masking")
	}

	// Test disabling all masking
	SetAllMasking(false)
	if defaultLogger.sensitiveMode != SHOW_SENSITIVE {
		t.Error("SetAllMasking(false) should disable sensitive masking")
	}
	if defaultLogger.piiMode != SHOW_PII {
		t.Error("SetAllMasking(false) should disable PII masking")
	}
}

// TestSetDevelopmentMode tests development mode configuration
func TestSetDevelopmentMode(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		level:           ERROR,
		writer:          io.Discard,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: []string{},
		piiFields:       []string{},
	}

	SetDevelopmentMode()

	if defaultLogger.format != PLAIN_FORMAT {
		t.Error("SetDevelopmentMode should set plain format")
	}
	if defaultLogger.sensitiveMode != SHOW_SENSITIVE {
		t.Error("SetDevelopmentMode should show sensitive data")
	}
	if defaultLogger.piiMode != SHOW_PII {
		t.Error("SetDevelopmentMode should show PII data")
	}
}

// TestSetProductionMode tests production mode configuration
func TestSetProductionMode(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		level:           DEBUG,
		writer:          io.Discard,
		format:          PLAIN_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE,
		piiMode:         SHOW_PII,
		sensitiveFields: []string{},
		piiFields:       []string{},
	}

	SetProductionMode()

	if defaultLogger.format != JSON_FORMAT {
		t.Error("SetProductionMode should set JSON format")
	}
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Error("SetProductionMode should mask sensitive data")
	}
	if defaultLogger.piiMode != MASK_PII {
		t.Error("SetProductionMode should mask PII data")
	}
}

// TestTimestampPrecision tests timestamp precision configuration
func TestTimestampPrecision(t *testing.T) {
	// Test ParseTimestampPrecision
	testCases := []struct {
		input    string
		expected TimestampPrecision
	}{
		{"second", SecondPrecision},
		{"sec", SecondPrecision},
		{"millisecond", MillisecondPrecision},
		{"milli", MillisecondPrecision},
		{"microsecond", MicrosecondPrecision},
		{"micro", MicrosecondPrecision},
		{"nanosecond", NanosecondPrecision},
		{"nano", NanosecondPrecision},
		{"invalid", NanosecondPrecision}, // default to highest precision
	}

	for _, tc := range testCases {
		result := ParseTimestampPrecision(tc.input)
		if result != tc.expected {
			t.Errorf("ParseTimestampPrecision(%s) = %v, want %v", tc.input, result, tc.expected)
		}
	}

	// Test SetTimestampPrecisionConfig and GetTimestampPrecisionConfig
	SetTimestampPrecisionConfig(NanosecondPrecision)
	if GetTimestampPrecisionConfig() != NanosecondPrecision {
		t.Error("SetTimestampPrecisionConfig should update timestamp precision")
	}

	SetTimestampPrecisionConfig(SecondPrecision)
	if GetTimestampPrecisionConfig() != SecondPrecision {
		t.Error("SetTimestampPrecisionConfig should update timestamp precision")
	}
}

// TestParseLogLevel tests log level parsing
func TestParseLogLevel(t *testing.T) {
	testCases := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", DEBUG},
		{"DEBUG", DEBUG},
		{"info", INFO},
		{"INFO", INFO},
		{"warn", WARN},
		{"warning", WARN},
		{"WARN", WARN},
		{"error", ERROR},
		{"ERROR", ERROR},
		{"invalid", INFO}, // default
	}

	for _, tc := range testCases {
		result := ParseLogLevel(tc.input)
		if result != tc.expected {
			t.Errorf("ParseLogLevel(%s) = %v, want %v", tc.input, result, tc.expected)
		}
	}
}

// TestClearFieldCache tests the ClearFieldCache function
func TestClearFieldCache(t *testing.T) {
	// Just call it to ensure it doesn't panic
	ClearFieldCache()
}

// TestSetUltraFastTimestampPrecision tests the SetUltraFastTimestampPrecision function
func TestSetUltraFastTimestampPrecision(t *testing.T) {
	// Test setting different interval values (in seconds)
	SetUltraFastTimestampPrecision(1)
	SetUltraFastTimestampPrecision(5)
	SetUltraFastTimestampPrecision(10)
	// Test with value < 1 (should be clamped to 1)
	SetUltraFastTimestampPrecision(0)
}

// TestSensitiveModeInvalidValue tests SetSensitiveMode with invalid value
func TestSensitiveModeInvalidValue(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		sensitiveMode: SHOW_SENSITIVE,
	}

	// Test with invalid mode - should default to MASK_SENSITIVE
	SetSensitiveMode("invalid_mode")
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Error("Invalid mode should default to MASK_SENSITIVE")
	}
}

// TestPIIModeInvalidValue tests SetPIIMode with invalid value
func TestPIIModeInvalidValue(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	defaultLogger = &Logger{
		piiMode: SHOW_PII,
	}

	// Test with invalid mode - should default to MASK_PII
	SetPIIMode("invalid_mode")
	if defaultLogger.piiMode != MASK_PII {
		t.Error("Invalid mode should default to MASK_PII")
	}
}

// TestInitFromEnvironment tests environment variable initialization
func TestInitFromEnvironment(t *testing.T) {
	originalLogger := defaultLogger
	defer func() { defaultLogger = originalLogger }()

	// Create a fresh logger for testing
	defaultLogger = &Logger{
		format:        JSON_FORMAT,
		level:         INFO,
		sensitiveMode: MASK_SENSITIVE,
		piiMode:       MASK_PII,
	}

	// Test EMIT_FORMAT=plain
	os.Setenv("EMIT_FORMAT", "plain")
	defer os.Unsetenv("EMIT_FORMAT")
	initFromEnvironment()
	if defaultLogger.format != PLAIN_FORMAT {
		t.Error("EMIT_FORMAT=plain should set PLAIN_FORMAT")
	}

	// Reset and test EMIT_FORMAT=json
	defaultLogger.format = PLAIN_FORMAT
	os.Setenv("EMIT_FORMAT", "json")
	initFromEnvironment()
	if defaultLogger.format != JSON_FORMAT {
		t.Error("EMIT_FORMAT=json should set JSON_FORMAT")
	}

	// Test EMIT_FORMAT=invalid (should default to JSON)
	os.Setenv("EMIT_FORMAT", "invalid")
	initFromEnvironment()
	if defaultLogger.format != JSON_FORMAT {
		t.Error("EMIT_FORMAT=invalid should default to JSON_FORMAT")
	}

	// Test EMIT_LEVEL
	os.Setenv("EMIT_LEVEL", "debug")
	defer os.Unsetenv("EMIT_LEVEL")
	initFromEnvironment()
	if defaultLogger.level != DEBUG {
		t.Error("EMIT_LEVEL=debug should set DEBUG level")
	}

	// Test EMIT_SHOW_CALLER
	os.Setenv("EMIT_SHOW_CALLER", "true")
	defer os.Unsetenv("EMIT_SHOW_CALLER")
	initFromEnvironment()
	if !defaultLogger.showCaller {
		t.Error("EMIT_SHOW_CALLER=true should enable caller display")
	}

	// Test EMIT_MASK_SENSITIVE=false
	os.Setenv("EMIT_MASK_SENSITIVE", "false")
	defer os.Unsetenv("EMIT_MASK_SENSITIVE")
	initFromEnvironment()
	if defaultLogger.sensitiveMode != SHOW_SENSITIVE {
		t.Error("EMIT_MASK_SENSITIVE=false should set SHOW_SENSITIVE")
	}

	// Test EMIT_MASK_SENSITIVE=true
	os.Setenv("EMIT_MASK_SENSITIVE", "true")
	initFromEnvironment()
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Error("EMIT_MASK_SENSITIVE=true should set MASK_SENSITIVE")
	}

	// Test EMIT_MASK_SENSITIVE=invalid (should default to MASK)
	os.Setenv("EMIT_MASK_SENSITIVE", "invalid")
	initFromEnvironment()
	if defaultLogger.sensitiveMode != MASK_SENSITIVE {
		t.Error("EMIT_MASK_SENSITIVE=invalid should default to MASK_SENSITIVE")
	}

	// Test EMIT_MASK_PII=false
	os.Setenv("EMIT_MASK_PII", "false")
	defer os.Unsetenv("EMIT_MASK_PII")
	initFromEnvironment()
	if defaultLogger.piiMode != SHOW_PII {
		t.Error("EMIT_MASK_PII=false should set SHOW_PII")
	}

	// Test EMIT_MASK_PII=true
	os.Setenv("EMIT_MASK_PII", "true")
	initFromEnvironment()
	if defaultLogger.piiMode != MASK_PII {
		t.Error("EMIT_MASK_PII=true should set MASK_PII")
	}

	// Test EMIT_MASK_PII=invalid (should default to MASK)
	os.Setenv("EMIT_MASK_PII", "invalid")
	initFromEnvironment()
	if defaultLogger.piiMode != MASK_PII {
		t.Error("EMIT_MASK_PII=invalid should default to MASK_PII")
	}

	// Test EMIT_MASK_STRING
	os.Setenv("EMIT_MASK_STRING", "***CUSTOM***")
	defer os.Unsetenv("EMIT_MASK_STRING")
	initFromEnvironment()
	if defaultLogger.maskString != "***CUSTOM***" {
		t.Error("EMIT_MASK_STRING should set custom mask string")
	}

	// Test EMIT_PII_MASK_STRING
	os.Setenv("EMIT_PII_MASK_STRING", "***PII_CUSTOM***")
	defer os.Unsetenv("EMIT_PII_MASK_STRING")
	initFromEnvironment()
	if defaultLogger.piiMaskString != "***PII_CUSTOM***" {
		t.Error("EMIT_PII_MASK_STRING should set custom PII mask string")
	}

	// Test EMIT_TIMESTAMP_PRECISION
	os.Setenv("EMIT_TIMESTAMP_PRECISION", "millisecond")
	defer os.Unsetenv("EMIT_TIMESTAMP_PRECISION")
	initFromEnvironment()
}
