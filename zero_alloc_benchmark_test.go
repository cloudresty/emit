package emit

import (
	"io"
	"testing"
	"time"
)

// Benchmarks for Zero-Allocation API (InfoZ, ErrorZ, etc.)

func BenchmarkInfoZ(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		InfoZ("User action completed")
	}
}

func BenchmarkInfoZWithFields(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard

	userID := "user123"
	count := 42
	active := true

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		InfoZ("User action completed",
			ZString("user_id", userID),
			ZInt("count", count),
			ZBool("active", active),
		)
	}
}

func BenchmarkInfoZWithSensitiveFields(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard

	email := "user@example.com"
	password := "secret123"
	sessionToken := "abc123def456"
	userID := "user123"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		InfoZ("User login",
			ZString("email", email),        // PII - should be masked
			ZString("password", password),  // Sensitive - should be masked
			ZString("token", sessionToken), // Sensitive - should be masked
			ZString("user_id", userID),     // Regular field
		)
	}
}

func BenchmarkInfoZWithManyFields(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		InfoZ("Complex operation",
			ZString("operation", "database_query"),
			ZString("table", "users"),
			ZInt("rows_affected", 123),
			ZInt64("execution_time_ns", int64(time.Millisecond)),
			ZFloat64("cpu_usage", 0.85),
			ZBool("success", true),
			ZTime("timestamp", time.Now()),
			ZDuration("latency", 50*time.Millisecond),
			ZString("query_id", "query_abc123"),
			ZString("session_id", "session_def456"),
		)
	}
}

// Comparison benchmarks against existing APIs

func BenchmarkInfoF_vs_InfoZ_Simple(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard

	b.Run("InfoF", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoF("User action completed", F().String("user_id", "user123"))
		}
	})

	b.Run("InfoZ", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("User action completed", ZString("user_id", "user123"))
		}
	})
}

func BenchmarkInfoF_vs_InfoZ_Complex(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard

	b.Run("InfoF", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoF("User action",
				F().String("user_id", "user123").
					Int("count", 42).
					Bool("active", true).
					String("email", "user@example.com"))
		}
	})

	b.Run("InfoZ", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("User action",
				ZString("user_id", "user123"),
				ZInt("count", 42),
				ZBool("active", true),
				ZString("email", "user@example.com"))
		}
	})
}

// Benchmarks for different output formats

func BenchmarkInfoZ_JSON_vs_Plain(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard

	b.Run("JSON", func(b *testing.B) {
		SetFormat("json")
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Test message",
				ZString("key1", "value1"),
				ZInt("key2", 42))
		}
	})

	b.Run("Plain", func(b *testing.B) {
		SetFormat("plain")
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Test message",
				ZString("key1", "value1"),
				ZInt("key2", 42))
		}
	})
}

// Benchmark against industry standard (Zap-style performance target)
func BenchmarkInfoZ_Target_Performance(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard
	SetFormat("json")

	// This benchmark targets Zap-level performance:
	// - Basic logging: 150-300 ns/op, <100 B/op, <3 allocs/op
	// - Structured logging: 300-600 ns/op, <200 B/op, <5 allocs/op

	b.Run("Basic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Basic message")
		}
	})

	b.Run("Structured", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Structured message",
				ZString("service", "api"),
				ZInt("status", 200),
				ZBool("cached", false))
		}
	})

	b.Run("Complex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Complex message",
				ZString("service", "api"),
				ZString("endpoint", "/users/profile"),
				ZInt("status_code", 200),
				ZInt64("response_time_ns", 1500000),
				ZFloat64("cpu_usage", 0.75),
				ZBool("cached", false),
				ZTime("request_time", time.Now()))
		}
	})
}

// Memory allocation pattern benchmarks
func BenchmarkInfoZ_AllocationPattern(b *testing.B) {
	SetLevel("info")
	defaultLogger.writer = io.Discard
	SetFormat("json")

	// Test different allocation patterns to optimize memory usage
	b.Run("NoFields", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Message with no fields")
		}
	})

	b.Run("SingleField", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Message with one field", ZString("key", "value"))
		}
	})

	b.Run("ThreeFields", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Message with three fields",
				ZString("key1", "value1"),
				ZInt("key2", 42),
				ZBool("key3", true))
		}
	})

	b.Run("FiveFields", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			InfoZ("Message with five fields",
				ZString("key1", "value1"),
				ZInt("key2", 42),
				ZBool("key3", true),
				ZFloat64("key4", 3.14),
				ZTime("key5", time.Now()))
		}
	})
}
