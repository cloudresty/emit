package main

import (
	"testing"
	"time"

	"github.com/cloudresty/emit"
)

// EmitBenchmarkSet contains all Emit-specific benchmarks
type EmitBenchmarkSet struct{}

// GetBenchmarks returns all Emit benchmark functions
func (e EmitBenchmarkSet) GetBenchmarks() []BenchmarkFunc {
	return []BenchmarkFunc{
		// Simple message benchmarks
		{"Emit_SimpleMessage", e.BenchmarkSimpleMessage},

		// StructuredFields benchmarks (zero-allocation)
		{"Emit_StructuredFields", e.BenchmarkStructuredFields},
		{"Emit_StructuredFieldsWithData", e.BenchmarkStructuredFieldsFields},
		{"Emit_StructuredFieldsComplex", e.BenchmarkStructuredFieldsComplex},

		// Emit-specific field benchmarks
		{"Emit_Field", e.BenchmarkField},
		{"Emit_FieldComplex", e.BenchmarkFieldComplex},

		// Key-value benchmarks
		{"Emit_KeyValue", e.BenchmarkKeyValue},
		{"Emit_KeyValueComplex", e.BenchmarkKeyValueComplex},

		// Memory-pooled benchmarks
		{"Emit_Pool", e.BenchmarkPool},
		{"Emit_PoolComplex", e.BenchmarkPoolComplex},

		// Security benchmarks
		{"Emit_SecurityBuiltIn", e.BenchmarkSecurityBuiltIn},
		{"Emit_SecurityDisabled", e.BenchmarkSecurityDisabled},
	}
}

// Simple message logging benchmarks
func (e EmitBenchmarkSet) BenchmarkSimpleMessage(b *testing.B) {

	for b.Loop() {
		emit.Info.Msg("Simple log message")
	}
}

// StructuredFields benchmarks
func (e EmitBenchmarkSet) BenchmarkStructuredFields(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.StructuredFields("User action",
			emit.ZString("user_id", "12345"),
			emit.ZString("action", "login"),
			emit.ZString("ip_address", "192.168.1.100"),
			emit.ZBool("success", true))
	}
}

func (e EmitBenchmarkSet) BenchmarkStructuredFieldsFields(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.StructuredFields("User action",
			emit.ZString("user_id", "12345"),
			emit.ZString("action", "login"),
			emit.ZString("ip_address", "192.168.1.100"),
			emit.ZBool("success", true))
	}
}

func (e EmitBenchmarkSet) BenchmarkStructuredFieldsComplex(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.StructuredFields("Complex operation",
			emit.ZString("service", "user-service"),
			emit.ZString("operation", "create_user"),
			emit.ZString("user_id", "12345"),
			emit.ZString("email", "user@example.com"),
			emit.ZString("ip_address", "192.168.1.100"),
			emit.ZInt("status_code", 201),
			emit.ZFloat64("duration_ms", 15.75),
			emit.ZBool("success", true),
			emit.ZTime("timestamp", time.Now()),
			emit.ZString("correlation_id", "corr_abc123"))
	}
}

// Emit-specific field benchmarks
func (e EmitBenchmarkSet) BenchmarkField(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.Field("User action",
			emit.NewFields().
				String("user_id", "12345").
				String("action", "login").
				String("ip_address", "192.168.1.100").
				Bool("success", true))
	}
}

func (e EmitBenchmarkSet) BenchmarkFieldComplex(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.Field("Complex operation",
			emit.NewFields().
				String("service", "user-service").
				String("operation", "create_user").
				String("user_id", "12345").
				String("email", "user@example.com").
				String("ip_address", "192.168.1.100").
				Int("status_code", 201).
				Float64("duration_ms", 15.75).
				Bool("success", true).
				Time("timestamp", time.Now()).
				String("correlation_id", "corr_abc123"))
	}
}

// Key-value benchmarks
func (e EmitBenchmarkSet) BenchmarkKeyValue(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.KeyValue("User action",
			"user_id", "12345",
			"action", "login",
			"ip_address", "192.168.1.100",
			"success", true)
	}
}

func (e EmitBenchmarkSet) BenchmarkKeyValueComplex(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.KeyValue("Complex operation",
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

// Memory-pooled benchmarks
func (e EmitBenchmarkSet) BenchmarkPool(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.Pool("User action", func(pf *emit.PooledFields) {
			pf.String("user_id", "12345").
				String("action", "login").
				String("ip_address", "192.168.1.100").
				Bool("success", true)
		})
	}
}

func (e EmitBenchmarkSet) BenchmarkPoolComplex(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.Pool("Complex operation", func(pf *emit.PooledFields) {
			pf.String("service", "user-service").
				String("operation", "create_user").
				String("user_id", "12345").
				String("email", "user@example.com").
				String("ip_address", "192.168.1.100").
				Int("status_code", 201).
				Float64("duration_ms", 15.75).
				Bool("success", true).
				Time("timestamp", time.Now()).
				String("correlation_id", "corr_abc123")
		})
	}
}

// Security benchmarks - This is where the real difference shows!
func (e EmitBenchmarkSet) BenchmarkSecurityBuiltIn(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.StructuredFields("User registration",
			emit.ZString("password", "super_secret_123"), // Automatically masked (sensitive)
			emit.ZString("email", "user@example.com"),    // Automatically masked (PII)
			emit.ZString("api_key", "sk_live_abc123"),    // Automatically masked (sensitive)
			emit.ZString("user_id", "12345"),
			emit.ZString("ssn", "123-45-6789"), // Automatically masked (PII)
			emit.ZString("action", "register"),
		)
	}
}

func (e EmitBenchmarkSet) BenchmarkSecurityDisabled(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		emit.Info.StructuredFields("User registration",
			emit.ZString("user_password", "super_secret_123"), // NOT masked (different key name)
			emit.ZString("user_mail", "user@example.com"),     // NOT masked (different key name)
			emit.ZString("service_key", "sk_live_abc123"),     // NOT masked (different key name)
			emit.ZString("user_id", "12345"),
			emit.ZString("social_sec", "123-45-6789"), // NOT masked (different key name)
			emit.ZString("action", "register"),
		)
	}
}
