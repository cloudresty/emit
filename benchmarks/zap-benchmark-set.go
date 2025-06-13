package main

import (
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

// ZapBenchmarkSet contains all Zap-specific benchmarks
type ZapBenchmarkSet struct{}

// GetBenchmarks returns all Zap benchmark functions
func (z ZapBenchmarkSet) GetBenchmarks() []BenchmarkFunc {
	return []BenchmarkFunc{
		// Simple message benchmarks
		{"Zap_SimpleMessage", z.BenchmarkSimpleMessage},

		// Sugared logger benchmarks (closest to zero-alloc)
		{"Zap_SugaredLogger", z.BenchmarkSugaredLogger},
		{"Zap_SugaredLoggerFields", z.BenchmarkSugaredLoggerFields},
		{"Zap_SugaredLoggerFieldsComplex", z.BenchmarkSugaredLoggerFieldsComplex},

		// Structured field benchmarks
		{"Zap_StructuredFields", z.BenchmarkStructuredFields},
		{"Zap_StructuredFieldsComplex", z.BenchmarkStructuredFieldsComplex},

		// Security benchmarks
		{"Zap_SecurityNone", z.BenchmarkSecurityNone},
		{"Zap_SecurityManual", z.BenchmarkSecurityManual},
	}
}

// Simple message logging benchmarks
func (z ZapBenchmarkSet) BenchmarkSimpleMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		zapLogger.Info("Simple log message")
	}
}

// Sugared logger benchmarks (closest equivalent to Emit's zero-alloc)
func (z ZapBenchmarkSet) BenchmarkSugaredLogger(b *testing.B) {
	sugar := zapLogger.Sugar()
	b.ResetTimer()
	for b.Loop() {
		sugar.Info("Zero allocation log message")
	}
}

func (z ZapBenchmarkSet) BenchmarkSugaredLoggerFields(b *testing.B) {
	sugar := zapLogger.Sugar()
	b.ResetTimer()
	for b.Loop() {
		sugar.Infow("User action",
			"user_id", "12345",
			"action", "login",
			"ip_address", "192.168.1.100",
			"success", true)
	}
}

func (z ZapBenchmarkSet) BenchmarkSugaredLoggerFieldsComplex(b *testing.B) {
	sugar := zapLogger.Sugar()
	b.ResetTimer()
	for b.Loop() {
		sugar.Infow("Complex operation",
			"service", "user-service",
			"operation", "create_user",
			"user_id", "12345",
			"email", "user@example.com",
			"ip_address", "192.168.1.100",
			"status_code", 201,
			"duration_ms", 15.75,
			"success", true,
			"timestamp", time.Now(),
			"correlation_id", "corr_abc123")
	}
}

// Structured field benchmarks (standard Zap API)
func (z ZapBenchmarkSet) BenchmarkStructuredFields(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zapLogger.Info("User action",
			zap.String("user_id", "12345"),
			zap.String("action", "login"),
			zap.String("ip_address", "192.168.1.100"),
			zap.Bool("success", true))
	}
}

func (z ZapBenchmarkSet) BenchmarkStructuredFieldsComplex(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zapLogger.Info("Complex operation",
			zap.String("service", "user-service"),
			zap.String("operation", "create_user"),
			zap.String("user_id", "12345"),
			zap.String("email", "user@example.com"),
			zap.String("ip_address", "192.168.1.100"),
			zap.Int("status_code", 201),
			zap.Float64("duration_ms", 15.75),
			zap.Bool("success", true),
			zap.Time("timestamp", time.Now()),
			zap.String("correlation_id", "corr_abc123"))
	}
}

// Security benchmarks
func (z ZapBenchmarkSet) BenchmarkSecurityNone(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		zapLogger.Info("User registration",
			zap.String("password", "super_secret_123"), // EXPOSED!
			zap.String("email", "user@example.com"),    // EXPOSED!
			zap.String("api_key", "sk_live_abc123"),    // EXPOSED!
			zap.String("user_id", "12345"),
			zap.String("ssn", "123-45-6789"), // EXPOSED!
			zap.String("action", "register"),
		)
	}
}

func (z ZapBenchmarkSet) BenchmarkSecurityManual(b *testing.B) {
	maskSensitive := func(s string) string {
		if len(s) <= 4 {
			return "***"
		}
		return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
	}

	b.ResetTimer()
	for b.Loop() {
		zapLogger.Info("User registration",
			zap.String("password", maskSensitive("super_secret_123")), // Manually masked
			zap.String("email", maskSensitive("user@example.com")),    // Manually masked
			zap.String("api_key", maskSensitive("sk_live_abc123")),    // Manually masked
			zap.String("user_id", "12345"),
			zap.String("ssn", maskSensitive("123-45-6789")), // Manually masked
			zap.String("action", "register"),
		)
	}
}
