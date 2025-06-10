package emit

import (
	"io"
	"testing"
)

// High-performance optimization benchmarks

// BenchmarkOptimizedMasking compares old vs new masking implementation
func BenchmarkOptimizedMasking(b *testing.B) {
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

	b.Run("Original", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = logger.maskSensitiveFields(fields)
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = logger.maskSensitiveFieldsFast(fields)
		}
	})
}

// BenchmarkOptimizedJSON compares old vs new JSON formatting
func BenchmarkOptimizedJSON(b *testing.B) {
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
		component:       "bench-test",
		version:         "v1.0.0",
	}

	fields := map[string]any{
		"user_id":    12345,
		"email":      "user@example.com",
		"password":   "secret123",
		"session_id": "sess_123456789",
		"safe_field": "safe_value",
	}

	b.Run("Original", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.logJSON(INFO, "Benchmark test message", fields)
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.logJSONFast(INFO, "Benchmark test message", fields)
		}
	})
}

// BenchmarkOptimizedPlain compares old vs new plain text formatting
func BenchmarkOptimizedPlain(b *testing.B) {
	logger := &Logger{
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

	fields := map[string]any{
		"user_id":    12345,
		"email":      "user@example.com",
		"password":   "secret123",
		"session_id": "sess_123456789",
		"safe_field": "safe_value",
	}

	b.Run("Original", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.logPlain(INFO, "Benchmark test message", fields)
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.logPlainFast(INFO, "Benchmark test message", fields)
		}
	})
}

// BenchmarkPooledFields tests the pooled fields performance
func BenchmarkPooledFields(b *testing.B) {
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

	b.Run("Traditional", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoWithFields("User operation", map[string]any{
				"user_id":  "12345",
				"email":    "user@example.com",
				"password": "secret123",
				"status":   200,
				"active":   true,
			})
		}
	})

	b.Run("FieldsBuilder", func(b *testing.B) {
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
	})

	b.Run("PooledFields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoFP("User operation", func(pf *PooledFields) {
				pf.String("user_id", "12345").
					String("email", "user@example.com").
					String("password", "secret123").
					Int("status", 200).
					Bool("active", true)
			})
		}
	})
}

// BenchmarkEndToEndOptimized tests full optimized pipeline performance
func BenchmarkEndToEndOptimized(b *testing.B) {
	// Create optimized logger that uses all fast paths
	optimizedLogger := &Logger{
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

	// Override log method to use optimized formatters
	originalLogger := defaultLogger
	defaultLogger = optimizedLogger
	defer func() { defaultLogger = originalLogger }()

	fields := map[string]any{
		"user_id":    12345,
		"email":      "user@example.com",
		"password":   "secret123",
		"session_id": "sess_123456789",
		"safe_field": "safe_value",
	}

	b.Run("OriginalPipeline", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoWithFields("Benchmark test with fields", fields)
		}
	})

	b.Run("OptimizedPipeline", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Use optimized formatting directly
			optimizedLogger.logJSONFast(INFO, "Benchmark test with fields", fields)
		}
	})

	b.Run("PooledFieldsPipeline", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoFP("Benchmark test with fields", func(pf *PooledFields) {
				pf.Int("user_id", 12345).
					String("email", "user@example.com").
					String("password", "secret123").
					String("session_id", "sess_123456789").
					String("safe_field", "safe_value")
			})
		}
	})
}

// BenchmarkMemoryUsage focuses on allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
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

	// Test different field sizes
	smallFields := map[string]any{"key": "value"}
	mediumFields := map[string]any{
		"user_id": 123, "email": "test@example.com", "status": "active",
	}
	largeFields := map[string]any{
		"user_id": 123, "email": "test@example.com", "password": "secret",
		"session": "sess_123", "ip": "192.168.1.1", "agent": "browser",
		"timestamp": "2025-06-09", "action": "login", "success": true,
	}

	b.Run("SmallFields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoWithFields("Small fields test", smallFields)
		}
	})

	b.Run("MediumFields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoWithFields("Medium fields test", mediumFields)
		}
	})

	b.Run("LargeFields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoWithFields("Large fields test", largeFields)
		}
	})
}
