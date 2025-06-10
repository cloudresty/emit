package emit

import (
	"io"
	"testing"
)

// BenchmarkInfoJSON benchmarks JSON logging without fields
func BenchmarkInfoJSON(b *testing.B) {
	// Create a logger that writes to discard (no I/O overhead)
	testLogger := &Logger{
		level:           INFO,
		writer:          io.Discard,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
		component:       "bench-test",
		version:         "v1.0.0",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Info("Benchmark test message")
	}
}

// BenchmarkInfoPlain benchmarks plain text logging without fields
func BenchmarkInfoPlain(b *testing.B) {
	testLogger := &Logger{
		level:           INFO,
		writer:          io.Discard,
		format:          PLAIN_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
		component:       "bench-test",
		version:         "v1.0.0",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Info("Benchmark test message")
	}
}

// BenchmarkInfoWithFieldsJSON benchmarks JSON logging with fields and masking
func BenchmarkInfoWithFieldsJSON(b *testing.B) {
	testLogger := &Logger{
		level:           INFO,
		writer:          io.Discard,
		format:          JSON_FORMAT,
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
		component:       "bench-test",
		version:         "v1.0.0",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	// Pre-create the fields map to avoid allocation overhead in benchmark
	fields := map[string]any{
		"user_id":    12345,
		"email":      "user@example.com",
		"password":   "secret123",
		"session_id": "sess_123456789",
		"safe_field": "safe_value",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		InfoWithFields("Benchmark test with fields", fields)
	}
}

// BenchmarkInfoWithFieldsNoMasking benchmarks JSON logging with fields but no masking
func BenchmarkInfoWithFieldsNoMasking(b *testing.B) {
	testLogger := &Logger{
		level:           INFO,
		writer:          io.Discard,
		format:          JSON_FORMAT,
		sensitiveMode:   SHOW_SENSITIVE, // No masking
		piiMode:         SHOW_PII,       // No masking
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		component:       "bench-test",
		version:         "v1.0.0",
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	fields := map[string]any{
		"user_id":    12345,
		"email":      "user@example.com",
		"password":   "secret123",
		"session_id": "sess_123456789",
		"safe_field": "safe_value",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		InfoWithFields("Benchmark test with fields", fields)
	}
}

// BenchmarkLogLevelFiltering benchmarks the performance when logs are filtered out
func BenchmarkLogLevelFiltering(b *testing.B) {
	testLogger := &Logger{
		level:           ERROR, // Only ERROR level logs will be processed
		writer:          io.Discard,
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

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// This should be filtered out and return early
		Info("This message should be filtered")
	}
}

// BenchmarkNestedFieldMasking benchmarks masking performance with nested objects
func BenchmarkNestedFieldMasking(b *testing.B) {
	testLogger := &Logger{
		level:           INFO,
		writer:          io.Discard,
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

	// Complex nested structure
	fields := map[string]any{
		"user": map[string]any{
			"email":    "user@example.com",
			"password": "secret123",
			"profile": map[string]any{
				"full_name": "John Doe",
				"phone":     "+1-555-123-4567",
				"address":   "123 Main St",
			},
		},
		"session": map[string]any{
			"token":      "session_token_123",
			"api_key":    "sk-1234567890",
			"expires_at": "2025-12-31T23:59:59Z",
		},
		"metadata": map[string]any{
			"request_id": "req_123456",
			"timestamp":  "2025-06-10T10:30:45Z",
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		InfoWithFields("Benchmark nested masking", fields)
	}
}

// BenchmarkParseLogLevel benchmarks log level parsing
func BenchmarkParseLogLevel(b *testing.B) {
	levels := []string{"debug", "info", "warn", "error", "DEBUG", "INFO", "WARN", "ERROR"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		level := levels[i%len(levels)]
		ParseLogLevel(level)
	}
}

// BenchmarkMaskSensitiveFields benchmarks the field masking function directly
func BenchmarkMaskSensitiveFields(b *testing.B) {
	testLogger := &Logger{
		sensitiveMode:   MASK_SENSITIVE,
		piiMode:         MASK_PII,
		sensitiveFields: defaultSensitiveFields,
		piiFields:       defaultPIIFields,
		maskString:      "***MASKED***",
		piiMaskString:   "***PII***",
	}

	fields := map[string]any{
		"user_id":    12345,
		"email":      "user@example.com",
		"password":   "secret123",
		"session_id": "sess_123456789",
		"safe_field": "safe_value",
		"phone":      "+1-555-123-4567",
		"api_key":    "sk-1234567890",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		testLogger.maskSensitiveFields(fields)
	}
}
