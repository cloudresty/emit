package emit

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

// Comprehensive benchmarks to identify optimization opportunities

// BenchmarkFieldsBuilderPerformance tests the Fields builder performance
func BenchmarkFieldsBuilderPerformance(b *testing.B) {
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
		InfoF("User operation", F().
			String("user_id", "12345").
			String("email", "user@example.com").
			String("password", "secret123").
			Int("status", 200).
			Bool("active", true))
	}
}

// BenchmarkKeyValuePerformance tests the KV API performance
func BenchmarkKeyValuePerformance(b *testing.B) {
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
		InfoKV("User operation",
			"user_id", "12345",
			"email", "user@example.com",
			"password", "secret123",
			"status", 200,
			"active", true)
	}
}

// BenchmarkMaskingOverhead isolates the masking performance cost
func BenchmarkMaskingOverhead(b *testing.B) {
	logger := &Logger{
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

	fields := map[string]any{
		"user_id":    12345,
		"email":      "user@example.com",
		"password":   "secret123",
		"session_id": "sess_123456789",
		"safe_field": "safe_value",
		"api_key":    "sk-secret-key",
		"phone":      "+1-555-123-4567",
		"address":    "123 Main St",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = logger.maskSensitiveFields(fields)
	}
}

// BenchmarkJSONMarshalingOverhead isolates JSON marshaling cost
func BenchmarkJSONMarshalingOverhead(b *testing.B) {
	logger := &Logger{
		level:         INFO,
		writer:        io.Discard,
		format:        JSON_FORMAT,
		sensitiveMode: SHOW_SENSITIVE, // No masking
		piiMode:       SHOW_PII,       // No masking
		component:     "bench-test",
		version:       "v1.0.0",
	}

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
		logger.logJSON(INFO, "Benchmark test message", fields)
	}
}

// BenchmarkLogLevelCheckingOverhead tests the cost of level checking
func BenchmarkLogLevelCheckingOverhead(b *testing.B) {
	testLogger := &Logger{
		level:  ERROR, // Set to ERROR so INFO logs are filtered out
		writer: io.Discard,
		format: JSON_FORMAT,
	}

	originalLogger := defaultLogger
	defaultLogger = testLogger
	defer func() { defaultLogger = originalLogger }()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Info("This message should be filtered out")
	}
}

// BenchmarkStringBuilding compares different string building approaches
func BenchmarkStringBuilding(b *testing.B) {
	fields := map[string]any{
		"user_id": 12345,
		"status":  "active",
		"count":   42,
	}

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for k, v := range fields {
				_ = fmt.Sprintf("%s=%v", k, v)
			}
		}
	})

	b.Run("strings.Builder", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for k, v := range fields {
				builder.WriteString(k)
				builder.WriteString("=")
				builder.WriteString(fmt.Sprintf("%v", v))
			}
			_ = builder.String()
		}
	})
}

// BenchmarkFieldsCreation compares different field creation methods
func BenchmarkFieldsCreation(b *testing.B) {
	b.Run("Traditional", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = map[string]any{
				"user_id": 12345,
				"email":   "user@example.com",
				"status":  200,
				"active":  true,
			}
		}
	})

	b.Run("FieldsBuilder", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = F().
				Int("user_id", 12345).
				String("email", "user@example.com").
				Int("status", 200).
				Bool("active", true).
				ToMap()
		}
	})

	b.Run("PreallocatedMap", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fields := make(map[string]any, 4) // Pre-allocated capacity
			fields["user_id"] = 12345
			fields["email"] = "user@example.com"
			fields["status"] = 200
			fields["active"] = true
			_ = fields
		}
	})
}
